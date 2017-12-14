package rule

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-kit/kit/endpoint"
)

type activateRequest struct {
	id string
}

type activateResponse struct{}

func (r activateResponse) StatusCode() int {
	return http.StatusNoContent
}

func activateEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(activateRequest)

		return activateResponse{}, svc.Activate(req.id)
	}
}

type deactivateRequest struct {
	id string
}

type deactivateResponse struct{}

func (r deactivateResponse) StatusCode() int {
	return http.StatusNoContent
}

func deactivateEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(deactivateRequest)

		return deactivateResponse{}, svc.Deactivate(req.id)
	}
}

type getRequest struct {
	id string
}

func getEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(getRequest)

		r, err := svc.GetByID(req.id)
		if err != nil {
			return nil, err
		}

		return &responseRule{rule: r}, nil
	}
}

type responseBucket struct {
	bucket Bucket
}

func (r *responseBucket) MarshalJSON() ([]byte, error) {
	ps := []responseParameter{}

	for k, v := range r.bucket.Parameters {
		ps = append(ps, responseParameter{
			key:   k,
			value: v,
		})
	}

	return json.Marshal(struct {
		Name       string              `json:"name"`
		Parameters []responseParameter `json:"parameters"`
		Percentage int                 `json:"percentage"`
	}{
		Name:       r.bucket.Name,
		Parameters: ps,
		Percentage: r.bucket.Percentage,
	})
}

func (r *responseBucket) UnmarshalJSON(raw []byte) error {
	v := struct {
		Name       string              `json:"name"`
		Parameters []responseParameter `json:"parameters"`
		Percentage int                 `json:"percentage"`
	}{}

	if err := json.Unmarshal(raw, &v); err != nil {
		return err
	}

	pm := Parameters{}

	for _, p := range v.Parameters {
		pm[p.key] = p.value
	}

	r.bucket = Bucket{
		Name:       v.Name,
		Parameters: pm,
		Percentage: v.Percentage,
	}

	return nil
}

type responseParameter struct {
	key   string
	value interface{}
}

func (r responseParameter) MarshalJSON() ([]byte, error) {
	v := struct {
		Name  string      `json:"name"`
		Type  string      `json:"type"`
		Value interface{} `json:"value"`
	}{
		Name:  r.key,
		Value: r.value,
	}

	switch r.value.(type) {
	case bool:
		v.Type = "bool"
	case float64:
		v.Type = "number"
	case string:
		v.Type = "string"
	default:
		v.Type = "unknown"
	}

	return json.Marshal(v)
}

func (r *responseParameter) UnmarshalJSON(raw []byte) error {
	v := struct {
		Name  string      `json:"name"`
		Value interface{} `json:"value"`
	}{}

	if err := json.Unmarshal(raw, &v); err != nil {
		return err
	}

	r.key = v.Name
	r.value = v.Value

	return nil
}

type responseRule struct {
	rule Rule
}

func (r *responseRule) MarshalJSON() ([]byte, error) {
	bs := []responseBucket{}

	for _, b := range r.rule.buckets {
		bs = append(bs, responseBucket{bucket: b})
	}

	return json.Marshal(struct {
		Active      bool             `json:"active"`
		ActivatedAt time.Time        `json:"activated_at"`
		Buckets     []responseBucket `json:"buckets"`
		ConfigID    string           `json:"config_id"`
		CreatedAt   string           `json:"created_at"`
		Criteria    *Criteria        `json:"criteria,omitempty"`
		Description string           `json:"description"`
		Deleted     bool             `json:"deleted"`
		EndTime     time.Time        `json:"end_time"`
		ID          string           `json:"id"`
		Kind        Kind             `json:"kind"`
		Name        string           `json:"name"`
		Rollout     uint8            `json:"rollout"`
		StartTime   time.Time        `json:"start_time"`
		UpdatedAt   time.Time        `json:"updated_at"`
	}{
		Active:      r.rule.active,
		ActivatedAt: r.rule.activatedAt,
		Buckets:     bs,
		ConfigID:    r.rule.configID,
		CreatedAt:   r.rule.createdAt.Format(time.RFC3339Nano),
		Criteria:    r.rule.criteria,
		Description: r.rule.description,
		Deleted:     r.rule.deleted,
		EndTime:     r.rule.endTime,
		ID:          r.rule.ID,
		Kind:        r.rule.kind,
		Name:        r.rule.name,
		Rollout:     r.rule.rollout,
		StartTime:   r.rule.startTime,
		UpdatedAt:   r.rule.updatedAt,
	})
}

func (r *responseRule) UnmarshalJSON(raw []byte) error {
	v := struct {
		Active      bool             `json:"active"`
		ActivatedAt time.Time        `json:"activated_at"`
		Buckets     []responseBucket `json:"buckets"`
		ConfigID    string           `json:"config_id"`
		CreatedAt   string           `json:"created_at"`
		Criteria    *Criteria        `json:"criteria,omitempty"`
		Description string           `json:"description"`
		Deleted     bool             `json:"deleted"`
		EndTime     time.Time        `json:"end_time"`
		ID          string           `json:"id"`
		Kind        Kind             `json:"kind"`
		Name        string           `json:"name"`
		Rollout     uint8            `json:"rollout"`
		StartTime   time.Time        `json:"start_time"`
		UpdatedAt   time.Time        `json:"updated_at"`
	}{}

	if err := json.Unmarshal(raw, &v); err != nil {
		return err
	}

	bs := []Bucket{}

	for _, rb := range v.Buckets {
		bs = append(bs, rb.bucket)
	}

	r.rule = Rule{
		active:      v.Active,
		activatedAt: v.ActivatedAt,
		buckets:     bs,
		configID:    v.ConfigID,
		criteria:    v.Criteria,
		description: v.Description,
		deleted:     v.Deleted,
		endTime:     v.EndTime,
		ID:          v.ID,
		kind:        v.Kind,
		name:        v.Name,
		rollout:     v.Rollout,
		startTime:   v.StartTime,
		updatedAt:   v.UpdatedAt,
	}

	return nil
}

type updateRolloutRequest struct {
	id      string
	rollout uint8
}

type updateRolloutResponse struct{}

func (r updateRolloutResponse) StatusCode() int {
	return http.StatusNoContent
}

func updateRolloutEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(updateRolloutRequest)

		return updateRolloutResponse{}, svc.UpdateRollout(req.id, req.rollout)
	}
}
