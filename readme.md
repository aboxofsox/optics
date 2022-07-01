# Optics 
A rewrite of [gostman](https://github.com/aboxofsox/gostman) with a more clear intention.

## Usage
```
Use Optics to test application endpoints directly from your terminal.

Usage:
  optics [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  init        Initialize optics.toml
  run         Run optics

Flags:
  -h, --help     help for optics
  -y, --skip     Use an empy config.
  -t, --toggle   Help message for toggle

Use "optics [command] --help" for more information about a command.
```

## Usage Explained
Use Optics to test your API endpoints without you having to setup an HTTP client, or run your application. A response of each endpoint path is saved as a JSON file in the directory defined as `outdir` in `optics.toml`, along with a log file containing the response result for each.

## Config Example
```toml
name = "pexels"
scheme = "https"
host = "api.pexels.com"
endpoints = [ "v1/search"]
outfile = true
outdir = "res"

[query_params]
query = "nature"

[headers]
Authorization = "{{.API_KEY}}"
```
I chose to use TOML for a few of reasons, mostly because I'm tired of JSON config files, I feel like TOML is more straight to the point, and I wanted a clear difference between the response output and the configuration.

| Key Name | Description | Data Type |
|------|------|------|
| name | Name of the endpoint. | `string` |
| scheme | Scheme of the endpoint. | `string` |
| host | The host URL. | `string` |
| endpoints | List of endpoint paths | `[]string` |
| query_params | Key/value pairs that represent URL query parameters. | `map[string]string` |
| headers | Key/value pairs that represent HTTP headers. | `map[string]string` |

## Environment Variables
By default, Optics will check for any `.env` files in your current working directory. If they exist, they will be loaded. Within the config file, they are defined as `{{.Key_Name}}`.

## Proxy
You can choose to use a proxy server to test your endpoints. This mostly for those who plan on using a reverse proxy in production. There are already inherent benefits of using a proxy server, but a lot of those benefits don't apply here since it's only used for testing purposes.

### Binaries
You can find executable binaries in [releases](https://github.com/aboxofsox/optics/releases). I included two build scripts, one in PowerShell, the other BASH. If there is a specific platform you need, you can add the platform to the list like `windows/386`, or bring your own build script.
#### Todo
- [ ] Better response info.
- [ ] Simple file parsing to extract URLs.
- [ ] Terminal table layout for response information.
- [x] Add option to use a proxy server.
