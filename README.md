# Paralus cli
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fparalus%2Fcli.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fparalus%2Fcli?ref=badge_shield)


CLI tool to interact with Paralus.

## Usage

[Download](https://github.com/paralus/cli/releases) the latest CLI binary and add it to your PATH.

`export PATH=$PATH:/usr/local/bin/pctl`

You can download the config using below command (RECOMMENDED)

```
pctl config download http://console.paralus.local
Enter Email: admin@paralus.local
Enter Password: 
CLI config stored at `/home/.paralus/cli/config.json`
```

OR 

Download the config from Paralus dashboard (*My Tools -> Download CLI Config*), this should be placed in the default directory under `$HOME/.paralus/cli`. *Create the directory if it doesn't exist.*

Sample configuration file looks like below:

```json
{
    "profile": "dev",
    "rest_endpoint": "console-ic-oss.dev.paralus-edge.net",
    "ops_endpoint": "console-ic-oss.dev.paralus-edge.net",
    "api_key": "9cfa2b7e009032dd1cd070fff811d59560a5ba28",
    "api_secret": "76f60059a2b6a97535da1394b57fe520c709e4c7f877ce4a4bd665924f6ced11",
    "project": "default",
    "organization": "exampleorg",
    "partner": "example"
}
```

## Features

The CLI allows you an additional way to interact with Paralus. There are commands that you can use to perform certain tasks, for example:

- `pctl create p project --desc "Description of the project"` to create a project
- `pctl create cluster imported sample-imported-cluster` to import a cluster
- `pctl get users` to list all users

For a complete list of commands, refer to our [CLI documentation](https://www.paralus.io/docs/usage/cli).

## Community & Support

- Visit [Paralus website](https://paralus.io) for the complete documentation and helpful links.
- Join our [Slack channel](https://join.slack.com/t/paralus/shared_invite/zt-1a9x6y729-ySmAq~I3tjclEG7nDoXB0A) to post your queries and discuss features.
- Tweet to [@paralus_](https://twitter.com/paralus_/) on Twitter.
- Create [GitHub Issues](https://github.com/paralus/cli/issues) to report bugs or request features.

## Contributing

The easiest way to start is to look at existing issues and see if there’s something there that you’d like to work on. You can filter issues with the label “Good first issue” which are relatively self sufficient issues and great for first time contributors.

Once you decide on an issue, please comment on it so that all of us know that you’re on it.

If you’re looking to add a new feature, raise a [new issue](https://github.com/paralus/cli/issues) and start a discussion with the community. Engage with the maintainers of the project and work your way through.


## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fparalus%2Fcli.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fparalus%2Fcli?ref=badge_large)