package log_monitoring

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/satyrius/gonx"
)

const logParserConfig = "$remote_addr $user_identifier $remote_user [$time_local] \"$request\" $status $bytes_sent"

func NewGonxParser() Parser {
	return &gonxParserImpl{
		gonx.NewParser(logParserConfig),
	}
}

type gonxParserImpl struct {
	parser *gonx.Parser
}

func (p *gonxParserImpl) Parse(line string) (*LogEntry, error) {
	entity, err := p.parser.ParseString(line)
	if err != nil {
		return nil, err
	}

	req, err := entity.Field("request")
	if err != nil {
		return nil, err
	}

	status, err := entity.Field("status")
	if err != nil {
		return nil, err
	}

	reqParts := strings.Fields(req)
	if len(reqParts) != 3 {
		return nil, fmt.Errorf("request is invalid: %s", req)
	}

	bytesStr, err := entity.Field("bytes_sent")
	if err != nil {
		return nil, err
	}

	bytes, err := strconv.Atoi(bytesStr)
	if err != nil {
		return nil, err
	}

	entry := &LogEntry{
		Method:  reqParts[0],
		Path:    reqParts[1],
		Proto:   reqParts[2],
		Status:  status,
		Section: "/",
		Bytes:   bytes,
	}

	pathSections := strings.Split(entry.Path, "/")
	if len(pathSections) > 0 {
		entry.Section = fmt.Sprintf("/%s", pathSections[0])
	}

	return entry, nil
}
