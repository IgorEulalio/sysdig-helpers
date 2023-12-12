import requests
import csv
import os
import argparse
import json
# Setting the environment variables
sysdig_url = os.environ.get('SYSDIG_URL')
token = os.environ.get('SYSDIG_TOKEN')

def parse_arguments():
    parser = argparse.ArgumentParser(description='Fetch Sysdig cluster data.')
    parser.add_argument('--limit', type=int, default=150, help='Limit the number of results')
    parser.add_argument('--filter', type=str, default='', help='Filter criteria')
    parser.add_argument('--connected', type=str, default='', help='Connected status filter')
    parser.add_argument('--output', type=str, default='clusters.csv', help='Output file name')

    args = parser.parse_args()
    return args

# Function to get the cluster data with named arguments for filter, limit, and connected
def get_cluster_data(sysdig_url, token, limit, filter, connected):
    url = f"{sysdig_url}/api/cloud/v2/dataSources/clusters?limit={limit}&filter={filter}&connected={connected}"
    headers = {"Authorization": f"Bearer {token}"}
    response = requests.get(url, headers=headers)
    
    if response.status_code == 200:
        return response.json()
    else:
        return None

# Function to extract the required information and write to CSV
def write_to_csv_with_agent_data(file_name, clusters, sysdig_url, token):
    """
    Writes cluster data along with agent data to a CSV file. 
    Skips calling the second API for clusters that are disconnected.
    Also, skips agent_version if all agents have 'Never Connected' status.
    """
    with open(file_name, 'w', newline='') as file:
        writer = csv.writer(file)
        writer.writerow(['name', 'node_count', 'agentConnected', 'nodes_connected', 'agent_status', 'agent_version', 'provider', 'environment'])

        for cluster in clusters:
            name = cluster.get('name', '')
            node_count = cluster.get('nodeCount', 0)
            agent_connected = cluster.get('agentConnected', False)
            # retrieve only 4 character of name
            environment = get_environment(name[3])
            provider = cluster.get('provider', '')

            if agent_connected:
                agent_data = get_agent_data(sysdig_url, token, name)
                if agent_data:
                    nodes_connected = agent_data['agentStats']['totalCount']
                    agent_details = agent_data['details']

                    if agent_details:
                        # Filter out agents with 'Never Connected' status
                        relevant_agents = [agent for agent in agent_details if agent['agentStatus'] != 'Never Connected']
                        
                        if relevant_agents:
                            # Use the first relevant agent's status and version
                            agent_status = relevant_agents[0]['agentStatus']
                            agent_version = relevant_agents[0]['agentVersion']
                        else:
                            # If all agents are 'Never Connected', set 'N/A'
                            agent_status = 'N/A'
                            agent_version = 'N/A'
                    else:
                        agent_status = 'N/A'
                        agent_version = 'N/A'
                else:
                    nodes_connected = 'N/A'
                    agent_status = 'N/A'
                    agent_version = 'N/A'
            else:
                # For disconnected clusters, set the agent fields to 'N/A'
                nodes_connected = 0
                agent_status = 'N/A'
                agent_version = 'N/A'

            writer.writerow([name, node_count, agent_connected, nodes_connected, agent_status, agent_version, provider, environment])

def get_agent_data(sysdig_url, token, cluster_name):
    """
    Calls the API for agent data for a given cluster name.
    """
    url = f"{sysdig_url}/api/cloud/v2/dataSources/agents?filter={cluster_name}"
    headers = {"Authorization": f"Bearer {token}"}
    response = requests.get(url, headers=headers)

    if response.status_code == 200:
        return response.json()
    else:
        return None

# Function to return environment based on environment single string
def get_environment(environment):
    if environment == 'd':
        return 'development'
    elif environment == 'p':
        return 'production'
    elif environment == 'i':
        return 'pre-production'
    else:
        return 'unknown'

# Main execution
if __name__ == "__main__":
    args = parse_arguments()
    clusters = get_cluster_data(sysdig_url, token, args.limit, args.filter, args.connected)
    if clusters:
        write_to_csv_with_agent_data(args.output, clusters, sysdig_url, token)
    else:
        print("Failed to retrieve data from the API.")
