# Callwrapper

Callwrapper is a package that helps you to wrap your function with a middleware.

## Features
- <b>Caching support</b>. Cache the result of your function. This will help to improve latency. Help to reduce the load on your external service.
- <b>Customize Caching Client</b>. You can customize the caching client. We have provide some example on ```pkg/cache```.
- <b>Singleflight support</b>. You can use singleflight to prevent multiple requests from hitting the same function.
- <b>Hook support</b>. You can add hook to your function. This will help you to add some extra logic to your function.
	- <b>Before Hook</b>. This hook will be called before the function is called.
	- <b>After Hook</b>. This hook will be called after the function is called.
	- <b>OnErrorLog</b>. This hook will be called when the function returns an error.
    - <b>OnWarnLog</b>. This hook will be called when the function returns a warning.

## Roadmap
- <b>Support Circuit Breaker</b>. You can use circuit breaker to prevent your service to call external service when the external service is down.
- <b>Support Telemetry</b>. Help to monitoring your external call like latency, error rate and QPS.

## How to use
This is an example of how to use callwrapper. You can find the example in ```example/```.