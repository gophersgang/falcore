package falcore

import (
	"container/list"
	"http"
	"log"
	"reflect"
)

// Pipelines have an upstream and downstream list of filters.
// A request is passed through the upstream items in order UNTIL
// a Response is returned.  Once a request is returned, it is passed
// through ALL ResponseFilters in the Downstream list, in order.
//
// If no response is generated by any Filters a default 404 response is
// returned.
//
// The RequestDoneCallback (if set) will be called after the request 
// has completed.  The finished request object will be passed to
// the FilterRequest method for inspection.  Changes to the request
// will have no effect and the return value is ignored.
//
// 
type Pipeline struct {
	Upstream            *list.List
	Downstream          *list.List
	RequestDoneCallback RequestFilter
}

func NewPipeline() (l *Pipeline) {
	l = new(Pipeline)
	l.Upstream = list.New()
	l.Downstream = list.New()
	return
}

// Pipelines are also RequestFilters... wacky eh?
// Be careful though because a Pipeline will always returns a 
// response so no Filters after a Pipeline filter will be run.
func (p *Pipeline) FilterRequest(req *Request) *http.Response {
	return p.execute(req)
}

func (p *Pipeline) execute(req *Request) (res *http.Response) {
	for e := p.Upstream.Front(); e != nil && res == nil; e = e.Next() {
		if router, ok := e.Value.(Router); ok {
			t := reflect.TypeOf(router)
			req.startPipelineStage(t.String()[1:])
			pipe := router.SelectPipeline(req)
			req.finishPipelineStage()
			if pipe != nil {
				t := reflect.TypeOf(pipe)
				req.startPipelineStage(t.String()[1:])
				res = pipe.FilterRequest(req)
				req.finishPipelineStage()
			}
		} else if filter, ok := e.Value.(RequestFilter); ok {
			t := reflect.TypeOf(filter)
			req.startPipelineStage(t.String())
			res = filter.FilterRequest(req)
			req.finishPipelineStage()
			if res != nil {
				break
			}
		} else {
			log.Printf("%v is not a RequestFilter\n", e.Value)
			break
		}
	}

	if res == nil {
		// Error: No response was generated
		res = SimpleResponse(req.HttpRequest, 404, nil, "Not found\n")
	}

	p.down(req, res)
	return
}

func (p *Pipeline) down(req *Request, res *http.Response) {
	for e := p.Downstream.Front(); e != nil; e = e.Next() {
		if filter, ok := e.Value.(ResponseFilter); ok {
			t := reflect.TypeOf(filter)
			req.startPipelineStage(t.String())
			filter.FilterResponse(req, res)
			req.finishPipelineStage()
		} else {
			// TODO
			break
		}
	}
}
