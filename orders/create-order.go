package orders

import (
    "net/http"
    "encoding/json"
    "github.com/satori/go.uuid"
    "github.com/aws/aws-sdk-go/service/sns"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/aws"
    "log"
    "os"
)

var orderAcceptedTopicArn string = "arn:aws:sns:us-west-2:110303772622:ocs-order_accepted"
var snsClient *sns.SNS
var logger *log.Logger

func init() {
    sess := session.Must(session.NewSession())
    snsClient = sns.New(sess)
    logger = log.New(os.Stdout, "[orders]", 0)
}

type req struct {
    Name string `json:"name"`
    Phone string `json:"phone"`
    Address1 string `json:"address1"`
    Address2 string `json:"address2"`
    PostalCode string `json:"postalCode"`
    Instructions string `json:"additionalInstructions"`
}

type order struct {
    ID string `json:"id"`
    *req
}

func createOrderHandler(w http.ResponseWriter, r *http.Request)  {
    logger.Println("Create order request received.")
    decoder := json.NewDecoder(r.Body)
    var request req
    err := decoder.Decode(&request)
    if err != nil {
        logger.Println("JSON deserialization error.", err.Error())
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    logger.Printf("Request: %v\n", request)

    defer r.Body.Close()

    order := &order{
        ID: uuid.NewV4().String(),
        req: &request,
    }

    err = fireOrderAccepted(order)
    if err != nil {
        logger.Println("Error publishing OrderAccepte event.", err.Error())
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-type", "applicaton/json")
    encoder := json.NewEncoder(w)
    encoder.Encode(order)
    logger.Printf("Responded with order: %v\n", order)
}

func fireOrderAccepted(order *order) error {
    msg, err := json.Marshal(order)
    if err != nil {
        return err
    }

    params := &sns.PublishInput{
        TopicArn: aws.String(orderAcceptedTopicArn),
        Message: aws.String(string(msg)),
    }

    resp, err := snsClient.Publish(params)
    if err != nil {
        return err
    }

    logger.Println("Published OrderAccepted event", resp.GoString())

    return nil
}