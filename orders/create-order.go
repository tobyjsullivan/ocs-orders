package orders

import (
    "net/http"
    "encoding/json"
    "github.com/satori/go.uuid"
    "github.com/aws/aws-sdk-go/service/sns"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/aws"
)

var orderAcceptedTopicArn string = "arn:aws:sns:us-west-2:110303772622:ocs-order_accepted"
var snsClient *sns.SNS

func init() {
    sess := session.Must(session.NewSession())
    snsClient = sns.New(sess)
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
    decoder := json.NewDecoder(r.Body)
    var request req
    err := decoder.Decode(&request)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    defer r.Body.Close()

    order := &order{
        ID: uuid.NewV4().String(),
        req: &request,
    }

    err = fireOrderAccepted(order)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-type", "applicaton/json")
    encoder := json.NewEncoder(w)
    encoder.Encode(order)
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

    println("Published OrderAccepted event", resp.GoString())

    return nil
}