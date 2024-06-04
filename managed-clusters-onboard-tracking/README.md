# Managed Clusters Onboard Tracking

This repository contains a Go application to track the onboarding status of managed clusters. The app can filter and export onboarding data based on specific criteria.

## Usage

To run the application, use the following command:

```sh
go run main.go --filter <filter_option> --output <output_file>
```

### Command Options

- `--filter <filter_option>`: This option specifies the filter criteria for the onboarding data. Use this parameter to filter cluster names, example if you have clustera, clusterb and clusterba, using this parameter as **clusterb** would only provide data from clusterb and clusterba.

- `--output <output_file>`: This option specifies the name of the output CSV file where the filtered data will be saved. The output file should have a `.csv` extension.

## Prerequisites

- Go 1.16 or higher

## Installation

1. Clone the repository:

    ```sh
    git clone https://github.com/IgorEulalio/sysdig-helpers.git
    cd sysdig-helpers/managed-clusters-onboard-tracking
    ```

2. Install dependencies:

    ```sh
    go mod tidy
    ```

### Example

```sh
go run main.go --filter myclustername --output myfile.csv
```

In this example, the application filters the onboarding data for SBR clusters and saves the result to `myfile.csv`.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

