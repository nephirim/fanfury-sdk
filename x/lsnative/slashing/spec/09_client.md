<!--
order: 9
-->

# CLI

A user can query and interact with the `slashing` module using the CLI.

## Query

The `query` commands allow users to query `slashing` state.

```sh
furyd query slashing --help
```

### params

The `params` command allows users to query genesis parameters for the slashing module.

```sh
furyd query slashing params [flags]
```

Example:

```sh
furyd query slashing params
```

Example Output:

```yml
downtime_jail_duration: 600s
min_signed_per_window: "0.500000000000000000"
signed_blocks_window: "100"
slash_fraction_double_sign: "0.050000000000000000"
slash_fraction_downtime: "0.010000000000000000"
```

### signing-info

The `signing-info` command allows users to query signing-info of the validator using consensus public key.

```sh
furyd query slashing signing-infos [flags]
```

Example:

```sh
furyd query slashing signing-info '{"@type":"/cosmos.crypto.ed25519.PubKey","key":"Auxs3865HpB/EfssYOzfqNhEJjzys6jD5B6tPgC8="}'

```

Example Output:

```yml
address: cosmosvalcons1nrqsld3aw6lh6t082frdqc84uwxn0t958c
index_offset: "2068"
jailed_until: "1970-01-01T00:00:00Z"
missed_blocks_counter: "0"
start_height: "0"
tombstoned: false
```

### signing-infos

The `signing-infos` command allows users to query signing infos of all validators.

```sh
furyd query slashing signing-infos [flags]
```

Example:

```sh
furyd query slashing signing-infos
```

Example Output:

```yml
info:
- address: cosmosvalcons1nrqsld3aw6lh6t082frdqc84uwxn0t958c
  index_offset: "2075"
  jailed_until: "1970-01-01T00:00:00Z"
  missed_blocks_counter: "0"
  start_height: "0"
  tombstoned: false
pagination:
  next_key: null
  total: "0"
```

## Transactions

The `tx` commands allow users to interact with the `slashing` module.

```bash
furyd tx slashing --help
```

### unjail

The `unjail` command allows users to unjail a validator previously jailed for downtime.

```bash
  furyd tx slashing unjail --from mykey [flags]
```

Example:

```bash
furyd tx slashing unjail --from mykey
```

## gRPC

A user can query the `slashing` module using gRPC endpoints.

### Params

The `Params` endpoint allows users to query the parameters of slashing module.

```sh
cosmos.slashing.v1beta1.Query/Params
```

Example:

```sh
grpcurl -plaintext localhost:9090 cosmos.slashing.v1beta1.Query/Params
```

Example Output:

```json
{
  "params": {
    "signedBlocksWindow": "100",
    "minSignedPerWindow": "NTAwMDAwMDAwMDAwMDAwMDAw",
    "downtimeJailDuration": "600s",
    "slashFractionDoubleSign": "NTAwMDAwMDAwMDAwMDAwMDA=",
    "slashFractionDowntime": "MTAwMDAwMDAwMDAwMDAwMDA="
  }
}
```

### SigningInfo

The SigningInfo queries the signing info of given cons address.

```sh
cosmos.slashing.v1beta1.Query/SigningInfo
```

Example:

```sh
grpcurl -plaintext -d '{"cons_address":"cosmosvalcons1nrqsld3aw6lh6t082frdqc84uwxn0t958c"}' localhost:9090 cosmos.slashing.v1beta1.Query/SigningInfo
```

Example Output:

```json
{
  "valSigningInfo": {
    "address": "cosmosvalcons1nrqsld3aw6lh6t082frdqc84uwxn0t958c",
    "indexOffset": "3493",
    "jailedUntil": "1970-01-01T00:00:00Z"
  }
}
```

### SigningInfos

The SigningInfos queries signing info of all validators.

```sh
cosmos.slashing.v1beta1.Query/SigningInfos
```

Example:

```sh
grpcurl -plaintext localhost:9090 cosmos.slashing.v1beta1.Query/SigningInfos
```

Example Output:

```json
{
  "info": [
    {
      "address": "cosmosvalcons1nrqslkwd3pz096lh6t082frdqc84uwxn0t958c",
      "indexOffset": "2467",
      "jailedUntil": "1970-01-01T00:00:00Z"
    }
  ],
  "pagination": {
    "total": "1"
  }
}
```

## REST

A user can query the `slashing` module using REST endpoints.

### Params

```sh
/cosmos/slashing/v1beta1/params
```

Example:

```sh
curl "localhost:1317/cosmos/slashing/v1beta1/params"
```

Example Output:

```json
{
  "params": {
    "signed_blocks_window": "100",
    "min_signed_per_window": "0.500000000000000000",
    "downtime_jail_duration": "600s",
    "slash_fraction_double_sign": "0.050000000000000000",
    "slash_fraction_downtime": "0.010000000000000000"
}
```

### signing_info

```sh
/cosmos/slashing/v1beta1/signing_infos/%s
```

Example:

```sh
curl "localhost:1317/cosmos/slashing/v1beta1/signing_infos/cosmosvalcons1nrqslkwd3pz096lh6t082frdqc84uwxn0t958c"
```

Example Output:

```json
{
  "val_signing_info": {
    "address": "cosmosvalcons1nrqslkwd3pz096lh6t082frdqc84uwxn0t958c",
    "start_height": "0",
    "index_offset": "4184",
    "jailed_until": "1970-01-01T00:00:00Z",
    "tombstoned": false,
    "missed_blocks_counter": "0"
  }
}
```

### signing_infos

```sh
/cosmos/slashing/v1beta1/signing_infos
```

Example:

```sh
curl "localhost:1317/cosmos/slashing/v1beta1/signing_infos
```

Example Output:

```json
{
  "info": [
    {
      "address": "cosmosvalcons1nrqslkwd3pz096lh6t082frdqc84uwxn0t958c",
      "start_height": "0",
      "index_offset": "4169",
      "jailed_until": "1970-01-01T00:00:00Z",
      "tombstoned": false,
      "missed_blocks_counter": "0"
    }
  ],
  "pagination": {
    "next_key": null,
    "total": "1"
  }
}
```
