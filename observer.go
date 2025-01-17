// SPDX-FileCopyrightText: 2025 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

// Package wrpzap provides a simple observer that logs information about the
// message being processed.  The logging is designed to be consistent with the
// fields used by the wrp-go package for ease of understanding.
//
// Because this package is entirely public interfaces, it is easy to extend or
// modify the logging behavior.
//
// The package is designed to work with the wrp-go package's Observer, Processor,
// and Modifier interfaces.
package wrpzap

import (
	"github.com/xmidt-org/wrp-go/v3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Observer logs information about the message being processed and sends the
// message to the next handler in the chain.
type Observer struct {
	Logger  *zap.Logger
	Level   zapcore.Level
	Message string
	Fields  []FieldOpt
}

// ObserveWRP logs information about the message being processed.
func (ob Observer) ObserveWRP(msg wrp.Message) {
	if ob.Logger == nil {
		return
	}

	fields := make([]zap.Field, 0, len(ob.Fields))
	for _, field := range ob.Fields {
		fields = append(fields, field(msg))
	}

	ob.Logger.Log(ob.Level, ob.Message, fields...)
}

// FieldOpt is a function that returns a zap.Field based on the message.
type FieldOpt func(wrp.Message) zap.Field

// LogMessageType logs the message type as a number.
func LogMessageType() FieldOpt {
	return LogMessageTypeAsNum()
}

// LogMessageTypeAsString logs the message type as a string.
func LogMessageTypeAsString() FieldOpt {
	return func(msg wrp.Message) zap.Field {
		return zap.Stringer(fMsgType, msg.Type)
	}
}

// LogMessageTypeAsNum logs the message type as a number.
func LogMessageTypeAsNum() FieldOpt {
	return func(msg wrp.Message) zap.Field {
		return zap.Int(fMsgType, int(msg.Type))
	}
}

// LogSource logs the source of the message.
func LogSource() FieldOpt {
	return func(msg wrp.Message) zap.Field {
		return zap.String(fSource, msg.Source)
	}
}

// LogDestination logs the destination of the message.
func LogDestination() FieldOpt {
	return func(msg wrp.Message) zap.Field {
		return zap.String(fDestination, msg.Destination)
	}
}

// LogTransactionUUID logs the transaction UUID of the message.
func LogTransactionUUID() FieldOpt {
	return func(msg wrp.Message) zap.Field {
		return zap.String(fTransactionUUID, msg.TransactionUUID)
	}
}

// LogContentType logs the content type of the message.
func LogContentType() FieldOpt {
	return func(msg wrp.Message) zap.Field {
		return zap.String(fContentType, msg.ContentType)
	}
}

// LogAccept logs the accept header of the message.
func LogAccept() FieldOpt {
	return func(msg wrp.Message) zap.Field {
		return zap.String(fAccept, msg.Accept)
	}
}

// LogStatus logs the status of the message.
func LogStatus() FieldOpt {
	return func(msg wrp.Message) zap.Field {
		return zap.Int64p(fStatus, msg.Status)
	}
}

// LogRequestDeliveryResponse logs the request delivery response of the message.
func LogRequestDeliveryResponse() FieldOpt {
	return func(msg wrp.Message) zap.Field {
		return zap.Int64p(fRequestDeliveryResponse, msg.RequestDeliveryResponse)
	}
}

// LogHeaders logs the headers of the message.
func LogHeaders() FieldOpt {
	return func(msg wrp.Message) zap.Field {
		return zap.Strings(fHeaders, msg.Headers)
	}
}

// LogMetadata logs the metadata of the message.
func LogMetadata() FieldOpt {
	return func(msg wrp.Message) zap.Field {
		return zap.Any(fMetadata, msg.Metadata)
	}
}

// LogPath logs the path of the message.
func LogPath() FieldOpt {
	return func(msg wrp.Message) zap.Field {
		return zap.String(fPath, msg.Path)
	}
}

// LogPayload logs the payload of the message.
func LogPayload() FieldOpt {
	return func(msg wrp.Message) zap.Field {
		return zap.Binary(fPayload, msg.Payload)
	}
}

// LogPayloadSize logs the size of the payload of the message.
func LogPayloadSize() FieldOpt {
	return func(msg wrp.Message) zap.Field {
		return zap.Int(fPayloadSize, len(msg.Payload))
	}
}

// LogServiceName logs the service name of the message.
func LogServiceName() FieldOpt {
	return func(msg wrp.Message) zap.Field {
		return zap.String(fServiceName, msg.ServiceName)
	}
}

// LogURL logs the URL of the message.
func LogURL() FieldOpt {
	return func(msg wrp.Message) zap.Field {
		return zap.String(fURL, msg.URL)
	}
}

// LogPartnerIDs logs the partner IDs of the message.
func LogPartnerIDs() FieldOpt {
	return func(msg wrp.Message) zap.Field {
		return zap.Strings(fPartnerIDs, msg.PartnerIDs)
	}
}

// LogSessionID logs the session ID of the message.
func LogSessionID() FieldOpt {
	return func(msg wrp.Message) zap.Field {
		return zap.String(fSessionID, msg.SessionID)
	}
}

// LogQualityOfService logs the quality of service of the message.
func LogQualityOfService() FieldOpt {
	return func(msg wrp.Message) zap.Field {
		return zap.Int(fQualityOfService, int(msg.QualityOfService))
	}
}
