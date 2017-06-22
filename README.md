# Cartridge Skeleton
The purpose of this repository is to define a base cartridge with empty job definitions and pipelines to allow developers to rapidly develop their cartridges.

## Stucture
A cartridge is broken down into the following sections:

 * infra
  * For infrastructure-related items
 * jenkins
  * For Jenkins-related items
 * src
  * For source control-related items

## Metadata
Each cartridge should contain a "metadata.cartridge" file that specifies the following metadata:

 * `CARTRIDGE_SDK_VERSION`
  * This defines the version of the Cartridge SDK that the cartridge conforms to
 
## Using this Repository
When developing a cartridge it is advisable to make a copy of this repository and remove all of the README.md files so that it serves as a basis for the new cartridge.

# Test framework configuration file (draft)

## Container configuration options

- `limit`: maximum number of containers that can be deployed
- `memory`: how much memory the containers can take

Example:

````yaml
containers:
  limit: 10
  memory: 512mb
````

## Tests

- `framework`: the framework used to run the tests, example JUnit, Codepcetp, etc.
- `repo`: the repo where the tests are stored
- `args`: some extra arguments needed to run the tests

Example:

````yaml
junit:
  - repo: ssh://github@someurl.git
````

## Options
Options related to how the framework will behave and execute different tasks

- `time`: maximum time the framework should wait before starting new containers
- `priority`: which tests the framework should prioritize
- `processor`: the percentage the process should wait before splitting tests
- `memory`: how much memory the tests should consume before deploying new containers

Example:

````yaml
time: 3m
priority: time_consuming
````

## Sample file

````yaml
---
configuration:
  containers:
    limit: 10
    memory: 512mb
  tests:
    junit:
      - repo: ssh://github@someurl.git
    codecetp:
      - repo: ssh://github@someurl.git
  options:
    time: 3m
    priority: time_consuming
    processor: 30
````


