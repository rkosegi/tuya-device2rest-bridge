# Tuya device to REST bridge

Simple service to expose Tuya smart devices (such as plugs) using REST API.

Currently, only `DpQuery` (10) and `DpQueryNew` (16) is implemented.

## Security

Security is hard, so I won't even pretend :innocent:.
To secure this thing, use your favorite reverse proxy, such as nginx.
