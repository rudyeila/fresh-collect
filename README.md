# fresh-collect
Fresh collect is a tool that helps plan your shopping list based on recipes. 

The tool takes a list of recipes as input, parses all the ingredients and merges them, before adding the ingredients and quantities to your shopping list. 

Currently supported are recipes from HelloFresh (hence the name), but this could be expanded in the future.

For the shopping list, [Bring!](https://www.getbring.com/en/home) is supported. This could also be extended in the future. 

# Example
![Example usage](./assets/usage.gif)

# Setup
At the root level of this repository, create a `.env` file. In it write your Bring! credentials in the following format:
```
BRING_USER="<email-address>"
BRING_PASSWORD="<your-password>"
```

The HelloFresh API requires an access token. This is also found in the same `.env` file. For now it is simply manually added, since I am still not sure where this token is retrieved from. 
When visiting a Hello Fresh recipe in the browser, the requests made to the HelloFresh backend contain this token, but I could not directly see any authentication requests. The token seems to be hardcoded. 

# Install
Make sure `$GOPATH/bin` is added to your path. To find where you `$GOPATH` is, use `go env GOPATH`.

Pull this repository and CD into `./cli` 

Run `go install`
Ensure `ls $GOPATH/bin` now contains an exe called `fresh-collect`. 

You should now be able to call the tool from the CMD using `fresh-collect`.

# Quick Start
To use this tool, find a couple of HelloFresh recipes you want to cook, for example:
* https://www.hellofresh.de/recipes/wurziges-dal-curry-mit-smashed-potatoes-64df2a4d614f75555c20edba
* https://www.hellofresh.de/recipes/rindergeschnetzeltes-in-rosmarinsosse-58343e5dd4d92c5781367e02

Note the IDs of the recipes (the string at the end of the URL):
* 64df2a4d614f75555c20edba
* 58343e5dd4d92c5781367e02

Run the command to parse and merge the ingredients and store the output in a file:

`fresh-collect parse --output=ingredients.json -r 64df2a4d614f75555c20edba -r 58343e5dd4d92c5781367e02`

This parses the ingredients across the recipes and combines the quantites of ingredients with the same, e.g. Potatoes, which are in both recipes.
The result is written out as json to the specified file `ingredients.json`

Run the `add` command to add all the ingredients into your Bring! shopping list:
`fresh-collect add --input=ingredients.json --list=HelloFresh`

In my case, I want to add the ingredients to a list called HelloFresh, which I had already created. Note, the list must already exist.

# Usage
```
NAME:
   fresh-collect - A new cli application

USAGE:
   fresh-collect [global options] command [command options]

DESCRIPTION:
   a tool that helps plan your shopping list based on recipes

AUTHOR:
   Rudy Ailabouni <eilabouni.rudy@gmail.com>

COMMANDS:
   parse, p  Parse recipes by IDs and get grouped ingredients
   add, a    add ingredients to shopping list
   help, h   Shows a list of commands or help for one command

```

This module provides a CLI client for interacting with the API. For now there are two operations that can be done:
1. Parsing and merging ingredients from HelloFresh recipes and writing them to an output file
```
NAME:
   fresh-collect - A new cli application

USAGE:
   fresh-collect [global options] command [command options]

DESCRIPTION:
   a tool that helps plan your shopping list based on recipes

AUTHOR:
   Rudy Ailabouni <eilabouni.rudy@gmail.com>

COMMANDS:
   parse, p  Parse recipes by IDs and get grouped ingredients
   add, a    add ingredients to shopping list
   help, h   Shows a list of commands or help for one command

```

2. Adding the ingredients from the file to a shopping list in the Bring app. 
```
NAME:
   fresh-collect add - add ingredients to shopping list

USAGE:
   fresh-collect add [command options]

OPTIONS:
   --input value, -i value  JSON file containing the parsed ingredients and their quantities. Output of Parse command
   --list value, -l value   The name of the shopping list the ingredients should be added to
   --help, -h               show help
```
