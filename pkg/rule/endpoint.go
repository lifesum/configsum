package rule

import (
	"context"
	"encoding/json"
	"net/http"
	"sort"
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
	ps := ResponseParameters{}

	for n, v := range r.bucket.Parameters {
		ps = append(ps, ResponseParameter{
			Name:  n,
			Value: v,
		})
	}

	sort.Sort(ps)

	return json.Marshal(struct {
		Name       string             `json:"name"`
		Parameters ResponseParameters `json:"parameters"`
		Percentage int                `json:"percentage"`
	}{
		Name:       r.bucket.Name,
		Parameters: ps,
		Percentage: r.bucket.Percentage,
	})
}

func (r *responseBucket) UnmarshalJSON(raw []byte) error {
	v := struct {
		Name       string             `json:"name"`
		Parameters ResponseParameters `json:"parameters"`
		Percentage int                `json:"percentage"`
	}{}

	if err := json.Unmarshal(raw, &v); err != nil {
		return err
	}

	pm := Parameters{}

	for _, p := range v.Parameters {
		pm[p.Name] = p.Value
	}

	r.bucket = Bucket{
		Name:       v.Name,
		Parameters: pm,
		Percentage: v.Percentage,
	}

	return nil
}

// ResponseParameter used to represent a parameter on the wire.
type ResponseParameter struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

// MarshalJSON to satisfy json.Marshaler and include the value type of the
// parameter for clients to make easy decisions when materialising it.
func (r ResponseParameter) MarshalJSON() ([]byte, error) {
	v := struct {
		Name  string      `json:"name"`
		Type  string      `json:"type"`
		Value interface{} `json:"value"`
	}{
		Name:  r.Name,
		Value: r.Value,
	}

	switch r.Value.(type) {
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

// ResponseParameters is a collection of ResponseParameter.
type ResponseParameters []ResponseParameter

func (r ResponseParameters) Len() int {
	return len(r)
}

func (r ResponseParameters) Less(i, j int) bool {
	return r[i].Name > r[j].Name
}

func (r ResponseParameters) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

type responseCriteria struct {
	criteria *Criteria
}

func (r *responseCriteria) MarshalJSON() ([]byte, error) {
	var (
		locale string
		u      *responseCriteriaUser
	)

	if r.criteria.Locale != nil {
		locale = string(*r.criteria.Locale)
	}

	if r.criteria.User != nil {
		u = &responseCriteriaUser{user: r.criteria.User}
	}

	return json.Marshal(struct {
		Locale string                `json:"locale,omitempty"`
		User   *responseCriteriaUser `json:"user,omitempty"`
	}{
		Locale: locale,
		User:   u,
	})
}

func (r *responseCriteria) UnmarshalJSON(raw []byte) error {
	if string(raw) == "" || string(raw) == "{}" {
		return nil
	}

	v := struct {
		Locale string                `json:"locale,omitempty"`
		User   *responseCriteriaUser `json:"user,omitempty"`
	}{}

	if err := json.Unmarshal(raw, &v); err != nil {
		return err
	}

	r.criteria = &Criteria{}

	if v.Locale != "" {
		locale := MatcherString(v.Locale)

		r.criteria.Locale = &locale
	}

	if v.User != nil && v.User.user != nil {
		r.criteria.User = v.User.user
	}

	return nil
}

type responseCriteriaUser struct {
	user *CriteriaUser
}

func (r *responseCriteriaUser) MarshalJSON() ([]byte, error) {
	var (
		is           []string
		subscription *MatcherInt
	)

	if r.user.ID != nil {
		is = []string(*r.user.ID)
	}

	if r.user.Subscription != nil {
		subscription = r.user.Subscription
	}

	return json.Marshal(struct {
		ID           []string    `json:"id,omitempty"`
		Subscription *MatcherInt `json:"subscription,omitempty"`
	}{
		ID:           is,
		Subscription: subscription,
	})
}

func (r *responseCriteriaUser) UnmarshalJSON(raw []byte) error {
	if string(raw) == "" || string(raw) == "{}" {
		return nil
	}

	v := struct {
		ID           *MatcherStringList `json:"id"`
		Subscription *MatcherInt        `json:"subscription"`
	}{}

	if err := json.Unmarshal(raw, &v); err != nil {
		return err
	}

	r.user = &CriteriaUser{}

	if v.ID != nil {
		r.user.ID = v.ID
	}

	if v.Subscription != nil {
		r.user.Subscription = v.Subscription
	}

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

	var c *responseCriteria

	if r.rule.criteria != nil {
		c = &responseCriteria{criteria: r.rule.criteria}
	}

	return json.Marshal(struct {
		Active      bool              `json:"active"`
		ActivatedAt time.Time         `json:"activated_at"`
		Buckets     []responseBucket  `json:"buckets"`
		ConfigID    string            `json:"config_id"`
		CreatedAt   string            `json:"created_at"`
		Criteria    *responseCriteria `json:"criteria,omitempty"`
		Description string            `json:"description"`
		Deleted     bool              `json:"deleted"`
		EndTime     time.Time         `json:"end_time"`
		ID          string            `json:"id"`
		Kind        Kind              `json:"kind"`
		Name        string            `json:"name"`
		Rollout     uint8             `json:"rollout"`
		StartTime   time.Time         `json:"start_time"`
		UpdatedAt   time.Time         `json:"updated_at"`
	}{
		Active:      r.rule.active,
		ActivatedAt: r.rule.activatedAt,
		Buckets:     bs,
		ConfigID:    r.rule.configID,
		CreatedAt:   r.rule.createdAt.Format(time.RFC3339Nano),
		Criteria:    c,
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
		Active      bool              `json:"active"`
		ActivatedAt time.Time         `json:"activated_at"`
		Buckets     []responseBucket  `json:"buckets"`
		ConfigID    string            `json:"config_id"`
		CreatedAt   string            `json:"created_at"`
		Criteria    *responseCriteria `json:"criteria,omitempty"`
		Description string            `json:"description"`
		Deleted     bool              `json:"deleted"`
		EndTime     time.Time         `json:"end_time"`
		ID          string            `json:"id"`
		Kind        Kind              `json:"kind"`
		Name        string            `json:"name"`
		Rollout     uint8             `json:"rollout"`
		StartTime   time.Time         `json:"start_time"`
		UpdatedAt   time.Time         `json:"updated_at"`
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
		criteria:    v.Criteria.criteria,
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

type responseList struct {
	rules []Rule
}

func (r *responseList) MarshalJSON() ([]byte, error) {
	rs := []responseRule{}

	for _, r := range r.rules {
		rs = append(rs, responseRule{
			rule: r,
		})
	}

	return json.Marshal(struct {
		Rules []responseRule `json:"rules"`
	}{
		Rules: rs,
	})
}

func (r *responseList) UnmarshalJSON(raw []byte) error {
	v := struct {
		Rules []responseRule `json:"rules"`
	}{}

	if err := json.Unmarshal(raw, &v); err != nil {
		return err
	}

	for _, rr := range v.Rules {
		r.rules = append(r.rules, rr.rule)
	}

	return nil
}

func (r responseList) StatusCode() int {
	if len(r.rules) == 0 {
		return http.StatusNoContent
	}

	return http.StatusOK
}

func listEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		rs, err := svc.List()
		if err != nil {
			return nil, err
		}

		return &responseList{rules: rs}, nil
	}
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
