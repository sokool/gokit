package docs

import (
	"errors"
	"net/http"
	"sort"
	"text/template"
)

type readme struct {
	doc *Doc
}

func (m *readme) markdown(e Endpoints) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		endpoints := e.All()
		sort.Sort(byPath{endpoints})

		t, err := template.New("one-page").Parse(tMD)
		//t, err := template.ParseFiles("internal/platform/web/server/docs/readme.md")
		if err != nil {
			info("%s", err)
			return
		}
		//
		w.WriteHeader(http.StatusOK)
		if err := t.Execute(w, endpoints); err != nil {
			info("%s", err)
			return
		}
	}
}

func (m *readme) apib(s Endpoints) http.HandlerFunc { return m.markdown(s) }

func (m *readme) html(s Endpoints) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		endpoints := s.All()
		sort.Sort(byPath{endpoints})

		t, err := template.
			New("one-page").
			Funcs(template.FuncMap{"dict": dict}).
			Parse(tHTML)
		//f := "internal/platform/web/server/docs/readme.html"
		//t, err := template.
		//	New(path.Base(f)).
		//	Funcs(template.FuncMap{"dict": dict}).
		//	ParseFiles(f)

		if err != nil {
			info("%s", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := t.Execute(w, map[string]interface{}{
			"Name":      "Documentation",
			"Active":    !m.doc.off,
			"Endpoints": endpoints}); err != nil {
			info("%s", err)
			return
		}
	}
}

func dict(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("invalid dict call")
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}

var tMD = `
FORMAT: 1A

# Developers Platform Service
This is description of DPS project.

## Actions
{{- range . }}
  {{- if .Usage }}
+ [{{ .Usage }}](#endpoint-{{ .ID }})
  {{- end }}
{{- end }}

{{- range . }}
{{- if or .Usage }}
{{- $requestHeaders := index (.Definition.Objects "header") "header" }}
{{- $requestBody := .Definition.Request.JSON 10 }}
{{- if eq "null" $requestBody }}
    {{- $requestBody := "" }}      
{{- end }}
 
{{- $responseBody := .Definition.Response.JSON 10 }}
{{- $parameters := index (.Definition.Objects "parameter") "parameter" }}
          
# <a name="endpoint-{{ .ID }}"></a>{{ .Usage }} [{{ .Method }} {{ .Path }}]
{{- if .Description }}
    {{ "\r" }}{{ .Description }}{{ "\n" }}
{{- end}}

{{- if or .Definition.Request .Definition.Headers }}
+ Request (application/json)
  {{- if .Definition.Parameters }}
    + Parameters
    {{""}} 
    {{- range $k, $field := $parameters.Vars }}
      {{- if not $field.Ignored }}
        {{- $action := "optional" }}
        {{- if $field.Required }}
          {{- $action := "required" }}
        {{- end }}
        + {{ $field.NameR "parameter" }}: {{ $field.Example }} ({{ $field.Type }}, {{ $action }}) - {{ $field.Description }}
      {{- end }} 
    {{- end}}
    {{""}} 
  {{- end }}
  {{- if not .Definition.Headers.Root.Ignored  }}
    + Headers
     {{- range $k, $field := $requestHeaders.Vars }}
         {{- if not $field.Ignored }}             
             {{ $field.NameR "header" }}: {{ $field.Example }}
         {{- end }}
     {{- end }}
  {{- end }} 
  {{- if .Definition.Request }}
    + Attributes ({{ .Definition.Request.Root.Type }})
    {{- range $k, $field := .Definition.Request.ByName }}
      {{- if not $field.IsRoot }}
      {{- if not $field.Ignored }}
        {{- $action := "optional" }}
          {{- if $field.Required }}
            {{- $action := "required" }}
          {{- end }}
        {{- if $field.Object }}
        elo
        {{- end }}  
        + {{ $field.NormalisedPath }} ({{ $field.Type }}, {{ $action }}) - {{ $field.Description }}
      {{- end }}
      {{- end }} 
    {{- end}}
    {{ "" }}
  {{- end}}  
  {{- if ne $requestBody "null" }}
    + Body
    
            {{ .Definition.Request.JSON 12 }}
        
  {{- end }}
{{- end }}

{{- if .Definition.Response }}

+ Response 200 (application/json)

    + Attributes ({{ .Definition.Response.Root.Type }})
    {{- range $k, $field := .Definition.Response.ByName }}
      {{- if not $field.Ignored }}
      {{- if not $field.IsRoot }}
        + {{ $field.NormalisedPath }} ({{ $field.Type }}) - {{ $field.Description }}
      {{- end }}  
      {{- end }} 
    {{- end}}

    + Body  
          
            {{ .Definition.Response.JSON 12 }}
            
{{- end }}
{{- end }}
{{- end }}
`

var tHTML = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Documentation</title>
    <link href="//netdna.bootstrapcdn.com/font-awesome/4.0.3/css/font-awesome.css"
          rel="stylesheet">
    <style>
        *, :after, :before {
            box-sizing: border-box;
        }

        body {
            font-family: Circular, Helvetica, Arial, sans-serif;
            font-size: 16px;
            line-height: 1.5;
        }

        h1 {
            border-bottom: 1px solid #dddddd;
            padding-top: 30px;
            padding-left: 20px;
        }

        h3 {
            text-transform: capitalize;
        }

        table.endpoints {
            width: 100%;
            border: 0px solid;
        }

        .header {
            font-size: 2em;
            font-weight: bold;
        }

        table.endpoints {
            font-size: 0.8em;
        }

        table.endpoints thead th:first-child {
            width: 15%;
        }

        table.endpoints thead th:nth-child(2) {
            /*width: 12%;*/
        }

        table.endpoints thead th:nth-child(4) {
            /*width: 25%;*/
        }

        table.endpoints thead th:nth-child(5) {
            /*width: 1.5%;*/
        }

        table.endpoints thead th:nth-child(6) {
            /*width: 1.5%;*/
        }

        table.endpoints thead th {
            padding: unset;
            margin: unset;
            border-bottom: 2px solid #4a4a4a;
            color: #969696;
            font-weight: lighter;
        }

        table.endpoints tbody td:first-child {
            white-space: nowrap;
            text-align: left;
            color: #0e7000;
            font-weight: bold;
        }

        table.endpoints tbody td:nth-child(2) {
            white-space: nowrap;
            text-align: left;
            color: #f59b23;
        }

        table.endpoints tbody td:nth-child(3) {
            text-align: left;
            color: black;
        }

        table.endpoints tbody {
            display: table-row-group;
            vertical-align: middle;
            border-color: inherit;
        }

        table.endpoints tbody tr {
            display: table-row;
            vertical-align: inherit;
            border-color: inherit;
        }

        table.endpoints tbody td {
            padding: 5px 10px 5px 0;
            border-bottom: 1px solid #e6e6e6;
        }

        table.endpoints tbody td a {
            color: #717372;
        }

        blockquote {
            font-size: 1em;
            width: 100%;
            margin: 10px 0px;
            font-style: italic;
            color: #555555;
            padding: 0.5em 3px 0.5em 40px;
            border-left: 8px solid #aaaa;
            line-height: 1.6;
            position: relative;
            background: #EDEDED;
        }

        blockquote::before {
            content: "\201C";
            color: #aaaa;
            font-size: 4em;
            position: absolute;
            left: 0px;
            top: -20px;
        }

        blockquote::after {
            content: '';
        }

        blockquote span {
            display: block;
            color: #333333;
            font-style: normal;
            font-weight: bold;
            margin-top: 1em;
        }

        div.endpoint {
            margin: 50px 10px;
        }


        .collapsible {
            background-color: #777;
            color: white;
            cursor: pointer;
            padding: 10px;
            width: 100%;
            border: none;
            text-align: left;
            outline: none;
            font-size: 15px;
            font-weight: bold;
            border-bottom: 1px solid gray;
        }

        .active, .collapsible:hover {
            background-color: #555;
        }

        .content {
            padding: 0 10px;
            max-height: 0;
            overflow: hidden;
            transition: max-height 0.2s ease-out;
            background-color: #f1f1f1;
        }

        .edit {
            cursor: pointer;
        }

        .stuff {
            position: relative;
            top: -20px;
        }

        pre.prettyprint {
            border: none !important;
            font-size: 12px;
        }

        textarea:focus {

        }


        input:focus {
            outline: none;
            /*background: lightsteelblue;*/
            /*border-radius: 20px;*/
            /*padding-left: 15px;*/
        }

        input:not([value]) {
            color: lightgray;
        }


        .field-form input {
            border: none;
            background: inherit;
            color: black;
            font-weight: bold;
            width: 100%;

            text-overflow: ellipsis;
            white-space: nowrap;
            overflow: hidden;
        }

        .endpoint-form hr {
            margin-top: 1px;
            border-style: solid;
        }

        .endpoint-form input:disabled {
            color: inherit;
        }

        .endpoint-form input {
            /*padding-left: 15px;*/
            border: none;
            font-size: 1.7em;
            font-weight: bold;
            width: 100%;
        }

        .endpoint-form input[name="usage"] {
            width: 50%;
        }

        .endpoint-form input[name="access"] {
            width: 25%;
        }

        .endpoint-form input[name="return"] {
            width: 24%;
        }

        .endpoint-form textarea, textarea:focus {
            outline: none;
            width: 99%;
            background: transparent;
            font-size: inherit;
            font-style: inherit;
            color: inherit;
            border: none;
            resize: none;
        }

        .method-GET {
            color: green;
        }

        .method-PUT {
            color: dodgerblue;
        }

        .method-POST {
            color: #ffab23;
        }

        .method-DELETE {
            color: red;
        }

        .method-PATCH {
            color: sienna;
        }

        .object {
            color: red;
            text-align: left !important;
        }

        a.collecting-on {
            color: green!important;
        }
        a.collecting-off {
            color: red!important;
        }
    </style>

    <script src="https://cdn.rawgit.com/google/code-prettify/master/loader/run_prettify.js"></script>
    <script>

        let form = '' +
            '<form class="field-form">' +
            '<input placeholder="field description..." type="text" name="description">' +
            '<input type="checkbox" name="ignore"> ignore' +
            '<input type="checkbox" name="required"> required' +
            '<input type="submit" name="cta" value="save">' +
            '</form>';


    </script>
</head>
<body>

{{- define "object-table"}}
    {{ if .Object }}
        {{ $endpoint := .endpoint }}
        {{ $root := .Object.Root }}
        <div class="contentX">
            <form class="field-form">
                <table class="endpoints">
                    <thead>
                    <tr>
                        <th>KEY</th>
                        <th>TYPE</th>
                        <th>DESCRIPTION</th>
                        <th>EXAMPLE</th>
                        <th>I</th>
                        <th>R</th>
                    </tr>

                    </thead>
                    <tbody></tbody>
                    {{ range $id, $field := .Object.ByName }}

                        <tr class="tr-form{{- with $field.Object}} object{{- end }}"
                            id="field-{{ $field.ID }}"
                            {{ if not $field.IsRoot }}data-belongs="{{ $endpoint }}_{{ $field.Belongs }}"{{ end }}
                            data-path="{{ $endpoint }}_{{ $field.Path }}"
                            title="{{ $field.Name }} last updated at {{ $field.LastUsed }}">
                            <td>
                                <span style="user-select:none; white-space: pre;{{- with $field.Object}}color:red;cursor:pointer{{- end }}">{{ $field.IndentedName }}</span>
                            </td>
                            <td>
                                {{- if $field.Array }}[]{{- end }}<input
                                        type="text"
                                        name="type"
                                        value="{{ $field.Type }}">

                            </td>
                            <td>
                                <input type="hidden" name="object"
                                       value="{{ $root.Name }}">
                                <input type="hidden" name="endpoint"
                                       value="{{ $endpoint }}">
                                <input type="hidden" name="field"
                                       value="{{ $field.ID }}">
                                <input placeholder="field description..."
                                       value="{{ $field.Description }}"
                                       type="text"
                                       name="description"
                                       autocomplete="off">
                            </td>
                            <td>
                                <input placeholder="example value..."
                                       value="{{ $field.Example }}"
                                       type="text"
                                       name="example"
                                       autocomplete="off">
                            </td>
                            <td>
                                <input
                                        type="checkbox"
                                        name="ignore"
                                        value="{{ $field.Ignored }}"
                                        {{ if $field.Ignored }} checked {{ end }}>
                            </td>
                            <td>
                                <input
                                        type="checkbox"
                                        name="required"
                                        value="{{ $field.Required }}"
                                        {{ if $field.Required }} checked {{ end }}>
                            </td>
                        </tr>

                    {{ end }}

                </table>
            </form>
            {{/*<pre class="prettyprint">{{ .Object.Root.Typee }}<br>{{ .Object.JSON 12 }}</pre>*/}}
        </div>
    {{ end }}
{{- end }}


<table class="endpoints">
    <thead>
    <tr>
        <th>METHOD</th>
        <th>ENDPOINT</th>
        <th>USAGE</th>
        <th>RETURN</th>
        <th>ROLES</th>
        <th>COLLECTING {{- if .Active }} <a href="/v2/docs/activate/off">OFF</a>
            {{- else }} <a href="/v2/docs/activate/on">ON</a>
            {{- end }}
        </th>
    </tr>

    </thead>
    <tbody></tbody>
    {{ range .Endpoints }}
        {{ $response := .Definition.Response.Root }}
        <tr>
            <td><span class="method-{{ .Method }}">{{ .Method }}</span></td>
            <td><a href="#endpoint-{{ .ID }}">{{ .Path }}</a></td>
            <td><span style="color: gray">#{{ .Visits }}:</span> {{ .Usage }}</td>
            <td>{{- if $response.Array }}[]{{- end }}{{ $response.Type }}</td>
            <td>{{ .Roles }}</td>
            <td>
                <a href="/v2/docs/{{ .ID }}/activate" class="collecting-{{- if .Eager }}on{{- else }}off{{- end}}">
                    {{- if .Eager }}on{{- else }}off{{- end}}
                </a>
            </td>
        </tr>
    {{end}}
</table>


{{ range .Endpoints }}
    {{ $endpoint := .ID }}
    {{ $response := .Definition.Response.Root }}

    <div class="endpoint" id="endpoint-{{.ID}}">
        <form class="endpoint-form">
            <input placeholder="endpoint name..."
                   value="{{ .Usage }}"
                   type="text"
                   name="usage" autocomplete="off">
            <input placeholder="roles..."
                   value="{{ .Roles }}"
                   type="text"
                   name="access" autocomplete="off">
            <input placeholder="response..."
                   value="{{- if $response.Array }}[]{{- end }}{{ $response.Type }}"
                   type="text"
                   name="return"
                   disabled>
            <input type="hidden" value="{{ .ID }}" name="field">
            <hr>
            <blockquote>
                <span class="method-{{ .Method }}">{{ .Method }} {{ .Path }}</span>
                <textarea placeholder="put endpoint description here..."
                          name="description">{{ .Description }}</textarea>
            </blockquote>
        </form>

        {{ if .Definition.Parameters }}
            <button class="collapsible">PARAMETERS</button>
            <div class="content">{{ template "object-table" dict "Object" .Definition.Parameters "endpoint" $endpoint }}</div>
        {{ end }}

        {{ if .Definition.Query }}
            <button class="collapsible">QUERY</button>
            <div class="content">{{ template "object-table" dict "Object" .Definition.Query "endpoint" $endpoint }}</div>
        {{ end }}

        {{ if .Definition.Headers }}
            <button class="collapsible">REQUEST HEADERS</button>
            <div class="content">{{ template "object-table" dict "Object" .Definition.Headers "endpoint" $endpoint }}</div>
        {{ end }}

        {{ if .Definition.Request }}
            <button class="collapsible">REQUEST</button>
            <div class="content">{{ template "object-table" dict "Object" .Definition.Request "endpoint" $endpoint }}</div>
        {{ end }}

        {{ if .Definition.Response }}
            <button class="collapsible">RESPONSE</button>
            <div class="content">{{ template "object-table" dict "Object" .Definition.Response "endpoint" $endpoint }}</div>
        {{ end }}

    </div>
{{end}}


</body>

<script>
    let coll = document.getElementsByClassName("collapsible");
    let i;

    for (i = 0; i < coll.length; i++) {
        coll[i].addEventListener("click", function () {
            this.classList.toggle("active");
            let content = this.nextElementSibling;
            if (content.style.maxHeight) {
                content.style.maxHeight = null;
            } else {
                content.style.maxHeight = content.scrollHeight + "px";
            }
        });
    }

    let edits = document.getElementsByClassName("edit");
    for (let i = 0, size = edits.length; i < size; i++) {
        edits[i].addEventListener("click", e => {
                let t = e.target.closest(".edit");
                let description = e.target.parentElement.getElementsByClassName("description")[0];

                description.innerHTML = form;
                let f = description.getElementsByTagName("form")[0];


                f.addEventListener("submit", submit => {
                    submit.preventDefault();
                    let xhr = new XMLHttpRequest();
                    xhr.open('PATCH', '/v2/docs/' + t.dataset.endpoint + "/" + t.dataset.type + "/" + t.dataset.field, true);
                    xhr.send(JSON.stringify(formToJSON(f)));
                    description.innerHTML = f["description"].value
                });
            }
        )
    }

    const isCheckbox = element => element.type === 'checkbox';
    const isMultiSelect = element => element.options && element.multiple;
    const formToJSON = elements => [].reduce.call(elements, (data, element) => {
        if (isCheckbox(element)) {
            data[element.name] = element.checked
        } else if (isMultiSelect(element)) {
            data[element.name] = getSelectValues(element);
        } else {
            data[element.name] = element.value;
        }

        return data;
    }, {});

    let ef = document.getElementsByClassName("endpoint-form");
    for (let i = 0, size = ef.length; i < size; i++) {
        ef[i].addEventListener("change", () => {
            let data = formToJSON(ef[i]);
            let xhr = new XMLHttpRequest();
            xhr.open('PATCH', '/v2/docs/' + data.field, true);
            xhr.send(JSON.stringify(data));
        });

    }

    let ff = document.getElementsByClassName("tr-form");
    for (let i = 0, size = ff.length; i < size; i++) {
        ff[i].addEventListener("change", () => {
            let elems = [];
            for (let key in ff[i].childNodes) {
                if (ff[i].childNodes[key].innerHTML === undefined) {
                    continue;
                }

                if (ff[i].childNodes[key].childNodes[1] === undefined) {
                    continue;
                }

                for (var k in ff[i].childNodes[key].childNodes) {
                    if (ff[i].childNodes[key].childNodes[k].tagName === "INPUT") {
                        elems.push(ff[i].childNodes[key].childNodes[k]);
                    }
                }
            }

            let body = formToJSON(elems);
            let xhr = new XMLHttpRequest();
            xhr.open('PATCH', '/v2/docs/' + body.endpoint + "/" + body.object + "/" + body.field, true);
            xhr.send(JSON.stringify(body));


        });
    }

    let oo = document.getElementsByClassName("object");
    for (let i = 0, size = oo.length; i < size; i++) {
        oo[i].addEventListener("click", () => {
            let belongs = oo[i].dataset.path;
            if (oo[i].dataset.visible === undefined) oo[i].dataset.visible = "yes";
            let visible = oo[i].dataset.visible;
            let a = '[data-belongs*="' + belongs + '"]';

            if (visible === "yes") {
                oo[i].dataset.visible = "no";
            } else {
                oo[i].dataset.visible = "yes";
            }
            // &#8628;
            // oo[i].querySelector("span").innerHTML = "&#8625;";

            let fields = document.querySelectorAll(a);
            for (let i = 0, size = fields.length; i < size; i++) {
                if (visible === "yes") {
                    fields[i].style.display = "none";
                } else {
                    fields[i].style.display = "";
                }
            }

        });
    }

</script>
</html>
`
