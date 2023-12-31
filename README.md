# refill

Rehydrate JSON schemas from arbitrary text-based data using language models.

# Usage

```sh
$ go get github.com/joehewett/refill
```

```sh
$ refill -json example/bank.json -dir example/emls -instructions example/extra_info.txt -verbose
```

## Flags
- `-help` - Print help message
- `-json` - Path to JSON schema file containing empty JSON schema to be rehydrated.
- `-dir` - Path to directory containing files to process.
- `-file` - Path to an individual file to process.
- `-instructions` - Path to text file containing instructions for the LM to follow when rehydrating the schema. Useful for providing additional context to the LM about the data it is processing and the format of the desired output.
- `-verbose` - Print verbose output

## Environment Variables
- `OPENAI_API_KEY` - All users will need to set their OpenAI key, which can be found here: [https://platform.openai.com/api-keys](https://platform.openai.com/api-keys).
- `UNIDOC_LICENSE_API_KEY` - If you are parsing PDFs, you'll need to set your [UniDoc](https://unidoc.io/) license key. You can get a free license key [here](https://unidoc.io/).

# Example

Given this empty JSON schema containing PII:

```json
{
  "name": "",
  "email": "",
  "phone": "",
  "address": "",
}
```

And this text file:
```
Hi Steve, it's Joe from the office. I'm calling to let you know that your order is ready for pickup. Please come by the office at your earliest convenience to pick it up. Thanks!

J. Hewett
Sales Manager

joe@company.net
(555) 555-5555
```

The `refill` command will produce this output:

```json
{
  "name": "Joe Hewett",
  "email": "joe@company.net",
  "phone": "(555) 555-5555",
  "address": "",
}
```

# License

[MIT](LICENSE)
