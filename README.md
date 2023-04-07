# DNA Center

The Intent REST API that exposes specific capabilities of the Cisco DNA Center platform. 
The Intent API provides policy-based abstraction of business intent, allowing focus on an outcome rather than struggling with individual mechanisms steps.

![dna-center](images/dna-center.png)

## Useful Actions

### Know Your Domain

Use the Know Your Network REST APIs to GET information about clients, sites, topology, devices, and issues; Create (POST) and manage (PUT, DELETE) sites, devices, IP Pools, edge and border devices, and authentication profiles.

### Site Management

Use the Site Management APIs to do site design, adjust network settings, device onboarding, Software Image Management (SIM), Configuration Templates.

### Operational Tasks

Operational Tools enable access to manage backups, discover network devices, trace paths through the network, identify and discover network components using SNMP and other discovery protocols,


### Policy

The Cisco DNA Center platform enables creation of application policies used to reflect your organization's business intent and translates it into network-specific and device-specific configurations required by the different types, makes, models, operating systems, roles, and resource constraints of your network devices.

## Endpoints

Refer to the DNA Center connector [API specification](openapi.json) for details.

## Example Usage

< sdk example? >

## Authentication

Configuration requires your DNA Center Server URL, username and password. The connector automatically injects 'X-Auth-Token' header automatically for every requests made.

Reachability to the Server URL is required in order for the connector to work properly. Please configure your firewall rules appropriately.