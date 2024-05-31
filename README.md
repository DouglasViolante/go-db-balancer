
# go-db-balancer

Educational and skills honing project, oriented to a fictional scenario where two AWS DynamoDB tables have an unequal item count, therefore a balancing cronjob is started/trigged when pre-determined conditions are met, in order to prevent future performance issues or degration of services which uses the databases.

Logging Capability provided via ```log/slog``` library.
AWS Environment via ```LocalStack```.

## Third Party Libs

- [AWS SDK v2](https://aws.github.io/aws-sdk-go-v2/docs/)
- [Localstack](https://docs.localstack.cloud)

## Todo

- DockerFile
- Add the capability of balance Horizontal-Partitioned Tables.
- Increase the Total Number of Databases (which is currently two) supported.
- Add Telemetry.

## How to Run

- Requires Docker!

- Start LocalStack on Docker.

- Run the Script *prepare_dynamodb*, an useful script to create and populate DynamoDB tables.

- Run ```main.go```
