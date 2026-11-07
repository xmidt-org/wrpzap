// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package wrpzap

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xmidt-org/wrp-go/v3"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestObserver_ObserveWRP(t *testing.T) {
	num := int64(123)
	headers := []string{"header1", "header2"}
	payload := []byte("test payload")

	tests := []struct {
		name           string
		nilLogger      bool
		fields         []FieldOpt
		inputMessage   wrp.Message
		expectedFields []zap.Field
	}{
		{
			name:      "nil logger",
			nilLogger: true,
		}, {
			name:           "log message type as default (num)",
			fields:         []FieldOpt{LogMessageType()},
			inputMessage:   wrp.Message{Type: wrp.SimpleRequestResponseMessageType},
			expectedFields: []zap.Field{zap.Int(fMsgType, int(wrp.SimpleRequestResponseMessageType))},
		}, {
			name:           "log message type as num",
			fields:         []FieldOpt{LogMessageTypeAsNum()},
			inputMessage:   wrp.Message{Type: wrp.SimpleRequestResponseMessageType},
			expectedFields: []zap.Field{zap.Int(fMsgType, int(wrp.SimpleRequestResponseMessageType))},
		}, {
			name:           "log message type as string",
			fields:         []FieldOpt{LogMessageTypeAsString()},
			inputMessage:   wrp.Message{Type: wrp.SimpleRequestResponseMessageType},
			expectedFields: []zap.Field{zap.Stringer(fMsgType, wrp.SimpleRequestResponseMessageType)},
		}, {
			name:           "log source",
			fields:         []FieldOpt{LogSource()},
			inputMessage:   wrp.Message{Source: "test source"},
			expectedFields: []zap.Field{zap.String(fSource, "test source")},
		}, {
			name:           "log destination",
			fields:         []FieldOpt{LogDestination()},
			inputMessage:   wrp.Message{Destination: "test destination"},
			expectedFields: []zap.Field{zap.String(fDestination, "test destination")},
		}, {
			name:           "log transaction uuid",
			fields:         []FieldOpt{LogTransactionUUID()},
			inputMessage:   wrp.Message{TransactionUUID: "test uuid"},
			expectedFields: []zap.Field{zap.String(fTransactionUUID, "test uuid")},
		}, {
			name:           "log content type",
			fields:         []FieldOpt{LogContentType()},
			inputMessage:   wrp.Message{ContentType: "test content type"},
			expectedFields: []zap.Field{zap.String(fContentType, "test content type")},
		}, {
			name:           "log accept",
			fields:         []FieldOpt{LogAccept()},
			inputMessage:   wrp.Message{Accept: "test accept"},
			expectedFields: []zap.Field{zap.String(fAccept, "test accept")},
		}, {
			name:           "log nil status",
			fields:         []FieldOpt{LogStatus()},
			inputMessage:   wrp.Message{Status: nil},
			expectedFields: []zap.Field{zap.Int64p(fStatus, nil)},
		}, {
			name:           "log num status",
			fields:         []FieldOpt{LogStatus()},
			inputMessage:   wrp.Message{Status: &num},
			expectedFields: []zap.Field{zap.Int64p(fStatus, &num)},
		}, {
			name:           "log nil rdr",
			fields:         []FieldOpt{LogRequestDeliveryResponse()},
			inputMessage:   wrp.Message{RequestDeliveryResponse: nil},
			expectedFields: []zap.Field{zap.Int64p(fRequestDeliveryResponse, nil)},
		}, {
			name:           "log num rdr",
			fields:         []FieldOpt{LogRequestDeliveryResponse()},
			inputMessage:   wrp.Message{RequestDeliveryResponse: &num},
			expectedFields: []zap.Field{zap.Int64p(fRequestDeliveryResponse, &num)},
		}, {
			name:           "log headers",
			fields:         []FieldOpt{LogHeaders()},
			inputMessage:   wrp.Message{Headers: headers},
			expectedFields: []zap.Field{zap.Strings(fHeaders, headers)},
		}, {
			name:           "log metadata",
			fields:         []FieldOpt{LogMetadata()},
			inputMessage:   wrp.Message{Metadata: map[string]string{"key": "value"}},
			expectedFields: []zap.Field{zap.Any(fMetadata, map[string]string{"key": "value"})},
		}, {
			name:           "log path",
			fields:         []FieldOpt{LogPath()},
			inputMessage:   wrp.Message{Path: "test path"},
			expectedFields: []zap.Field{zap.String(fPath, "test path")},
		}, {
			name:           "log payload",
			fields:         []FieldOpt{LogPayload()},
			inputMessage:   wrp.Message{Payload: payload},
			expectedFields: []zap.Field{zap.Binary(fPayload, payload)},
		}, {
			name:           "log payload size",
			fields:         []FieldOpt{LogPayloadSize()},
			inputMessage:   wrp.Message{Payload: payload},
			expectedFields: []zap.Field{zap.Int(fPayloadSize, len(payload))},
		}, {
			name:           "log service name",
			fields:         []FieldOpt{LogServiceName()},
			inputMessage:   wrp.Message{ServiceName: "test service name"},
			expectedFields: []zap.Field{zap.String(fServiceName, "test service name")},
		}, {
			name:           "log url",
			fields:         []FieldOpt{LogURL()},
			inputMessage:   wrp.Message{URL: "test url"},
			expectedFields: []zap.Field{zap.String(fURL, "test url")},
		}, {
			name:           "log partner IDs",
			fields:         []FieldOpt{LogPartnerIDs()},
			inputMessage:   wrp.Message{PartnerIDs: []string{"partner1", "partner2"}},
			expectedFields: []zap.Field{zap.Strings(fPartnerIDs, []string{"partner1", "partner2"})},
		}, {
			name:           "log session ID",
			fields:         []FieldOpt{LogSessionID()},
			inputMessage:   wrp.Message{SessionID: "session123"},
			expectedFields: []zap.Field{zap.String(fSessionID, "session123")},
		}, {
			name:           "log quality of service",
			fields:         []FieldOpt{LogQualityOfService()},
			inputMessage:   wrp.Message{QualityOfService: 2},
			expectedFields: []zap.Field{zap.Int(fQualityOfService, 2)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, recorded := observer.New(zap.InfoLevel)
			logger := zap.New(core)

			text := "test message"
			ob := Observer{
				Logger:  logger,
				Level:   zap.InfoLevel,
				Message: text,
				Fields:  tt.fields,
			}

			if tt.nilLogger {
				ob.Logger = nil
			}

			ob.ObserveWRP(context.Background(), tt.inputMessage)

			entries := recorded.All()
			if tt.nilLogger {
				require.Empty(t, entries)
				return
			}
			require.Len(t, entries, 1)
			entry := entries[0]
			assert.Equal(t, zap.InfoLevel, entry.Level)
			assert.Equal(t, text, entry.Message)
			assert.ElementsMatch(t, tt.expectedFields, entry.Context)
		})
	}
}

func TestFieldOpt_JSONTags(t *testing.T) {
	fieldMap := map[string]string{
		"Type":                    fMsgType,
		"Source":                  fSource,
		"Destination":             fDestination,
		"TransactionUUID":         fTransactionUUID,
		"ContentType":             fContentType,
		"Accept":                  fAccept,
		"Status":                  fStatus,
		"RequestDeliveryResponse": fRequestDeliveryResponse,
		"Headers":                 fHeaders,
		"Metadata":                fMetadata,
		"Path":                    fPath,
		"Payload":                 fPayload,
		"ServiceName":             fServiceName,
		"URL":                     fURL,
		"PartnerIDs":              fPartnerIDs,
		"SessionID":               fSessionID,
		"QualityOfService":        fQualityOfService,
	}

	ignored := map[string]struct{}{
		"Spans":        {},
		"IncludeSpans": {},
	}

	msgType := reflect.TypeOf(wrp.Message{})

	for fieldName, expectedTag := range fieldMap {
		t.Run(fieldName, func(t *testing.T) {
			field, found := msgType.FieldByName(fieldName)
			require.True(t, found, "Field '%s' not found in wrp.Message", fieldName)

			jsonTag := field.Tag.Get("json")
			list := strings.SplitN(jsonTag, ",", 2)
			require.NotEmpty(t, list[0], "Field '%s' does not have a JSON tag", fieldName)

			assert.Equal(t, expectedTag, list[0], "Constant for field '%s' does not match the JSON tag", fieldName)
		})
	}

	// Ensure all fields in wrp.Message are represented in the fieldMap
	for i := range msgType.NumField() {
		field := msgType.Field(i)
		if _, found := ignored[field.Name]; found {
			continue
		}

		assert.Contains(t, fieldMap, field.Name, "Field '%s' is not represented in the fieldMap", field.Name)
	}
}
