# Golang Snapcast JSON RPC API

An implementation of the [Snapcast](https://github.com/badaix/snapcast) JSON RPC API. The Snapcast project is a multi-room audio player that allows synchronized playback across multiple devices.

For more information on the Snapcast JSON RPC API, refer to the [Snapcast JSON RPC documentation](https://github.com/badaix/snapcast/blob/develop/doc/json_rpc_api/control.md).

## Project Layout
- `snapcast/`
	> Types for api
- `snapclient/`
	> Full client implementation using [gorilla/websocket](https://github.com/gorilla/websocket)
