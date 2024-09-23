### Source

```puml
@startuml

interface mock.Router {
	Run(string) error
}

interface mock.Service {
    +Add(mock mockDTO) (*addMockResponse, error)
    +Match(request httpRequest) (*httpResponse, error)
}

interface mock.Repository {
+Save(mock.mock) error
+GetAll() []mock.mock
}

interface mock.Mocker {
   +When(req *mock.RequestDTO) mock.Expect
}

interface mock.Expect {
   +ThenReturn(resp *mock.ResponseDTO) error
}

class mock.MockDTO{
    ID: string
    Request : *RequestDTO
    Response : *ResponseDTO
}

class mock.RequestDTO{
    URL : map[string]string
    Method : *string
    Headers : map[string]map[string]string
    QueryParameters : map[string]map[string]string
    Priority : int
}

class mock.ResponseDTO{
    Status : int
    Body : json.RawMessage
    Headers : map[string]string
}

class mock.Mock{
+ID : string
+Request : RequestMatch
+Response : HttpResponse
}

class mock.HttpResponse{
    Status  :int
    Body    :[]byte
    Headers :map[string]string
}

class mock.RequestMatch {
    +URL             :SimpleCondition
    +Method          :*string
    +Headers         :ComplexConditions
    +QueryParameters :ComplexConditions
    +Body            :SimplexConditions
    +Priority        :int
    +IsExpected(request HttpRequest) bool
}

class mock.SimpleCondition{
    -operator:operator
    -value:string
}

class mock.ComplexCondition extends mock.SimpleCondition{
    -field:string
}

enum mock.Operator{
    undefined,
    equal,
    contains,
    pattern
}


mock.Router -down-> mock.Service :use
mock.Service -down-> mock.Repository : use
mock.Service -> mock.MockDTO : receive
mock.Service -> mock.Mock : makes
mock.Mocker -down-> mock.Expect : create
mock.Repository -> mock.Mock : store
mock.Mocker -> mock.RequestDTO : receive
mock.Expect -> mock.ResponseDTO : recive
mock.Expect -> mock.MockDTO : makes
mock.Expect -> mock.Service : use
mock.MockDTO::Request o-- mock.RequestDTO : has
mock.MockDTO::Response o-- mock.ResponseDTO : has
mock.SimpleCondition o-- mock.Operator : has
mock.RequestMatch::Headers *-- mock.ComplexCondition : has
mock.RequestMatch::QueryParameters *-- mock.ComplexCondition : has
mock.RequestMatch::URL o-- mock.SimpleCondition : has
mock.RequestMatch::Body o-- mock.SimpleCondition : has
mock.Mock::Request o-- mock.RequestMatch : has
mock.Mock::Response o-- mock.HttpResponse : has

@enduml
```


![Class Diagram](http://www.plantuml.com/plantuml/png/dLJ1Kjim43tNNw5Z4XC-m0S62ft2J3DGXquF30TXBKJ5bbn93cGw_VTANbiecuHsweD7Mhttz7ObxQbrtBYkgRDCAWVcaLV06bqz7vUwytFsA_jGTcfYdP7gQSh066sotple9PYrhC33cV_oCo4c8TulNjnUJzbaneLO-5a9jjNAmX5gJP6muAvQJGpyxC0wjdAkBN4Sc3Wbj5qF9RJQIgVDnjl8btmDapuUVgFp_3EuixgUJDdTVLpSenMSWi5qykyhKC4Rc-4-w2oaXU1FhorKZgh4XK7brgkKu3eJj6mRnJ2lNYrMjwgvjMzE_28MzkgYO7WK_NHmmWesUpCrhA0iBiKjY7Oa3VBVoY-UqF3s3ZUvJ_PQW5jf4VIIu2Lm0SPka_RpltPWDZVSy0RSFr9kZDH6ket7IkM8zoGhDBzqt7KMiNtWenQ1-iDgTLpobmLOov-0-gOxxoc_lqrN5onMFdyh-p3gRs-8nw9V9bUMAAPBMyvPiUNSVyNT_SF64Sj_Svnw6jzZEDZqAvQoQMiurqf89xL251vnWagFEXyDnDLp7JHUeuX573teGpXCadhQnSLRC7onIlcLnGS68dQEoxHpF6XTa-AGlB4ySzs2uKwROXnWUCth3egzHQOvCNXreCIxkpnAg4KY1wfh47yzR8l0JWbub0h4KJ_rsV0QXvLMZajbSTPovy2e89MbBNikz8kQdsmtvgApaAMH2AEzis2xq9EjNkGX1YgGwtTW_TDw1fl527M6P-TnbG7k8AC-Y9hrGTXMEw5TZ1h3BGAGjWn4KY3DCys1O6BIiLoHLXHZWz3p-HwBAsufVZ2M4Z126EHYHNoyGH5C6a5HZCztDbBe5JvC9I_wVoH20ufcQHO7MQ557AQHIjAo4pKAfrMdhNd0PwUWHDVKVm00)