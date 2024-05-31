package db

import (
	"context"
	config "go-db-balancer/configs"
	"go-db-balancer/internal/pkg/domain"
	"log/slog"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Adapter struct {
	db 	*dynamodb.Client
	sl  slog.Logger
}

func NewAdapter(logger slog.Logger) (*Adapter, error){

	cfg, err := config.GetAWSConfig()
	
	if err != nil {
		return nil, err
	}

	dbClient := dynamodb.NewFromConfig(cfg)

	return &Adapter{db: dbClient, sl: logger}, nil
}

func (adp Adapter) GetTablesItemCount(ctx context.Context, tableNames []string) ([]domain.Table, error) {

	tablesInfo := []domain.Table{}

	for _, table := range tableNames{
		params := dynamodb.ScanInput{
			TableName: 	aws.String(table),
		}
	
		count, err := adp.db.Scan(ctx, &params)
		if err != nil {
			adp.sl.Error("Unable to Execute Get Item Count Operation on Table %s, error: %s", table, err)
		}

		newTableInfo := domain.Table{TableName: table, ItemCount: len(count.Items)}
		tablesInfo = append(tablesInfo, newTableInfo)
	}

	return tablesInfo, nil
}

func (adp Adapter) GetItems(ctx context.Context, totalItems int, sourceTableName string) ([]domain.Specs, error) {
	params := dynamodb.ScanInput{
		TableName: 	aws.String(sourceTableName),
		Limit:     	aws.Int32(int32(totalItems)),
	}

	dynamoDbScanOutput, err := adp.db.Scan(ctx, &params)

	if err != nil {
		adp.sl.Error("Unable to Execute Get Operation on Table %s, error: %s", sourceTableName, err)
	} else if dynamoDbScanOutput.LastEvaluatedKey != nil{
		adp.sl.Warn("TODO: Not All Required Items were Retrieved from DynamoDB!")
	}

	specifications := []domain.Specs{}

	for _, i := range dynamoDbScanOutput.Items{
		record := domain.Specs{}
		
		err = attributevalue.UnmarshalMap(i, &record)

		if err != nil {
			adp.sl.Error("Error when Unmarshalling Data: %s", err)
			return []domain.Specs{}, err
		}

		specifications = append(specifications, record)
	}

	return specifications, nil
}

func (adp Adapter) PutItems(ctx context.Context, recordsToMove []domain.Specs, putTargetTableName string) error {

	for _, record := range recordsToMove {

		dataRow, err := attributevalue.MarshalMap(record)
 		if err != nil {
			adp.sl.Error("Error when Marshalling Data: %s\n", err)
  			return err
 		}
		input := &dynamodb.PutItemInput{
			TableName: &putTargetTableName,	
			Item: dataRow,
		}

		_, err = adp.db.PutItem(ctx, input)
		if err != nil {
			adp.sl.Error("Unable to Execute Put Operation on Table %s, error: %s", putTargetTableName, err)
		}
	}

	return nil
}

func (adp Adapter) DeleteItems(ctx context.Context, recordsToDelete []domain.Specs, deleteTargetTableName string) error {

	for _, record := range recordsToDelete {

        input := &dynamodb.DeleteItemInput{
			Key: map[string]types.AttributeValue{
				"ID": &types.AttributeValueMemberN{Value: *aws.String(strconv.Itoa(record.ID))},
			},
            TableName: &deleteTargetTableName,
        }

        _, err := adp.db.DeleteItem(ctx, input)
        if err != nil {
			adp.sl.Error("Unable to Execute Delete Operation on Table %s, error: %s", deleteTargetTableName, err)
        }
    }

	return nil
}