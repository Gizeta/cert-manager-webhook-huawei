package huawei

import (
	"errors"
	"fmt"
	"strings"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	dns "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/dns/v2"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/dns/v2/model"
	hwregion "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/dns/v2/region"
)

type Client struct {
	client *dns.DnsClient
}

func NewClient(ak, sk, region string) *Client {
	credentials := basic.NewCredentialsBuilder().
		WithAk(ak).
		WithSk(sk).
		Build()
	dnsClient := dns.NewDnsClient(
		dns.DnsClientBuilder().
			WithRegion(hwregion.ValueOf(region)).
			WithCredential(credentials).
			Build())
	return &Client{dnsClient}
}

func (c *Client) AddDomainRecord(domain string, recordName string, recordType string, value string) error {
	zoneId, err := c.GetZoneID(domain)
	if err != nil {
		return err
	}

	if recordType == "TXT" && !strings.HasPrefix(value, "\"") {
		value = fmt.Sprintf("\"%s\"", value)
	}

	request := &model.CreateRecordSetRequest{}
	request.ZoneId = zoneId
	request.Body = &model.CreateRecordSetReq{
		Records: []string{
			value,
		},
		Type: recordType,
		Name: recordName,
	}

	_, err = c.client.CreateRecordSet(request)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) DeleteDomainRecord(domain string, recordName string, recordType string) error {
	zoneId, err := c.GetZoneID(domain)
	if err != nil {
		return err
	}

	listRequest := &model.ListRecordSetsByZoneRequest{
		ZoneId: zoneId,
		Type:   &recordType,
		Name:   &recordName,
	}
	listResponse, err := c.client.ListRecordSetsByZone(listRequest)
	if err != nil {
		return err
	}
	if len(*listResponse.Recordsets) == 0 {
		return errors.New("record not found")
	}
	if len(*listResponse.Recordsets) > 1 {
		return errors.New("too many records")
	}

	deleteRequest := &model.DeleteRecordSetRequest{
		ZoneId:      zoneId,
		RecordsetId: *(*listResponse.Recordsets)[0].Id,
	}
	_, err = c.client.DeleteRecordSet(deleteRequest)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) GetZoneID(domain string) (string, error) {
	zone, err := c.client.ListPublicZones(&model.ListPublicZonesRequest{
		Name: &domain,
	})
	if err != nil {
		return "", err
	}
	if len(*zone.Zones) == 0 {
		return "", errors.New("zone not found")
	}
	return *(*zone.Zones)[0].Id, nil
}
