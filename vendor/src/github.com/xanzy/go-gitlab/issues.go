//
// Copyright 2015, Sander van Harmelen
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package gitlab

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

// IssuesService handles communication with the issue related methods
// of the GitLab API.
//
// GitLab API docs: http://doc.gitlab.com/ce/api/issues.html
type IssuesService struct {
	client *Client
}

// Issue represents a GitLab issue.
//
// GitLab API docs: http://doc.gitlab.com/ce/api/issues.html
type Issue struct {
	ID          int      `json:"id"`
	IID         int      `json:"iid"`
	ProjectID   int      `json:"project_id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Labels      []string `json:"labels"`
	Milestone   struct {
		ID          int       `json:"id"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		DueDate     string    `json:"due_date"`
		State       string    `json:"state"`
		UpdatedAt   time.Time `json:"updated_at"`
		CreatedAt   time.Time `json:"created_at"`
	} `json:"milestone"`
	Assignee struct {
		ID        int       `json:"id"`
		Username  string    `json:"username"`
		Email     string    `json:"email"`
		Name      string    `json:"name"`
		State     string    `json:"state"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"assignee"`
	Author struct {
		ID        int       `json:"id"`
		Username  string    `json:"username"`
		Email     string    `json:"email"`
		Name      string    `json:"name"`
		State     string    `json:"state"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"author"`
	State     string    `json:"state"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
}

func (i Issue) String() string {
	return Stringify(i)
}

// ListIssuesOptions represents the available ListIssues() options.
//
// GitLab API docs: http://doc.gitlab.com/ce/api/issues.html#list-issues
type ListIssuesOptions struct {
	ListOptions
	State   string   `url:"state,omitempty"`
	Labels  []string `url:"labels,omitempty"`
	OrderBy string   `url:"order_by,omitempty"`
	Sort    string   `url:"sort,omitempty"`
}

// ListIssues gets all issues created by authenticated user. This function
// takes pagination parameters page and per_page to restrict the list of issues.
//
// GitLab API docs: http://doc.gitlab.com/ce/api/issues.html#list-issues
func (s *IssuesService) ListIssues(opt *ListIssuesOptions) ([]*Issue, *Response, error) {
	req, err := s.client.NewRequest("GET", "issues", opt)
	if err != nil {
		return nil, nil, err
	}

	var i []*Issue
	resp, err := s.client.Do(req, &i)
	if err != nil {
		return nil, resp, err
	}

	return i, resp, err
}

// ListProjectIssuesOptions represents the available ListProjectIssues() options.
//
// GitLab API docs: http://doc.gitlab.com/ce/api/issues.html#list-issues
type ListProjectIssuesOptions struct {
	ListOptions
	IID       int      `url:"iid,omitempty"`
	State     string   `url:"state,omitempty"`
	Labels    []string `url:"labels,omitempty"`
	Milestone string   `url:"milestone,omitempty"`
	OrderBy   string   `url:"order_by,omitempty"`
	Sort      string   `url:"sort,omitempty"`
}

// ListProjectIssues gets a list of project issues. This function accepts
// pagination parameters page and per_page to return the list of project issues.
//
// GitLab API docs: http://doc.gitlab.com/ce/api/issues.html#list-project-issues
func (s *IssuesService) ListProjectIssues(
	pid interface{},
	opt *ListProjectIssuesOptions) ([]*Issue, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/issues", url.QueryEscape(project))

	req, err := s.client.NewRequest("GET", u, opt)
	if err != nil {
		return nil, nil, err
	}

	var i []*Issue
	resp, err := s.client.Do(req, &i)
	if err != nil {
		return nil, resp, err
	}

	return i, resp, err
}

// GetIssue gets a single project issue.
//
// GitLab API docs: http://doc.gitlab.com/ce/api/issues.html#single-issues
func (s *IssuesService) GetIssue(pid interface{}, issue int) (*Issue, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/issues/%d", url.QueryEscape(project), issue)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	i := new(Issue)
	resp, err := s.client.Do(req, i)
	if err != nil {
		return nil, resp, err
	}

	return i, resp, err
}

// CreateIssueOptions represents the available CreateIssue() options.
//
// GitLab API docs: http://doc.gitlab.com/ce/api/issues.html#new-issues
type CreateIssueOptions struct {
	Title       string   `url:"title,omitempty"`
	Description string   `url:"description,omitempty"`
	AssigneeID  int      `url:"assignee_id,omitempty"`
	MilestoneID int      `url:"milestone_id,omitempty"`
	Labels      []string `url:"labels,omitempty"`
}

// CreateIssue creates a new project issue.
//
// GitLab API docs: http://doc.gitlab.com/ce/api/issues.html#new-issues
func (s *IssuesService) CreateIssue(
	pid interface{},
	opt *CreateIssueOptions) (*Issue, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/issues", url.QueryEscape(project))

	// This is needed to get a single, comma separated string
	opt.Labels = []string{strings.Join(opt.Labels, ",")}

	req, err := s.client.NewRequest("POST", u, opt)
	if err != nil {
		return nil, nil, err
	}

	i := new(Issue)
	resp, err := s.client.Do(req, i)
	if err != nil {
		return nil, resp, err
	}

	return i, resp, err
}

// UpdateIssueOptions represents the available UpdateIssue() options.
//
// GitLab API docs: http://doc.gitlab.com/ce/api/issues.html#edit-issues
type UpdateIssueOptions struct {
	Title       string   `url:"title,omitempty"`
	Description string   `url:"description,omitempty"`
	AssigneeID  int      `url:"assignee_id,omitempty"`
	MilestoneID int      `url:"milestone_id,omitempty"`
	Labels      []string `url:"labels,omitempty"`
	StateEvent  string   `url:"state_event,omitempty"`
}

// UpdateIssue updates an existing project issue. This function is also used
// to mark an issue as closed.
//
// GitLab API docs: http://doc.gitlab.com/ce/api/issues.html#edit-issues
func (s *IssuesService) UpdateIssue(
	pid interface{},
	issue int,
	opt *UpdateIssueOptions) (*Issue, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/issues/%d", url.QueryEscape(project), issue)

	// This is needed to get a single, comma separated string
	opt.Labels = []string{strings.Join(opt.Labels, ",")}

	req, err := s.client.NewRequest("PUT", u, opt)
	if err != nil {
		return nil, nil, err
	}

	i := new(Issue)
	resp, err := s.client.Do(req, i)
	if err != nil {
		return nil, resp, err
	}

	return i, resp, err
}
