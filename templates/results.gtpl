<!DOCTYPE html>
<html>
<body>
<h1>Search Results</h1>
<p>Showing results for <b>{{.OriginalQuery.Query}}</b></p>
<p>Go <a href="/">home</a></p>
<p>Showing {{.NumResults}} of {{.TotalResults}}</p>
    <table border="1">
        <tr>
            <th>510(k) Number</th>
            <th>Title</th>
            <th>Applicant</th>
            <th>Hit Details</th>
            <th>Submission Date</th>
            <th>Predicates</th>
        </tr>
    {{ range .SearchResults }}
    <tr>
        <td><a href="https://www.accessdata.fda.gov/scripts/cdrh/cfdocs/cfPMN/pmn.cfm?ID={{.id}}">{{ .id }}</a></td>
        <td>{{ .title }}</td>
        <td>{{ .applicant }}</td>
        <td>{{unescapeHTML ._formatted.full_text}}</td>
        <td>{{ .submission_date }}</td>
        <td>{{ range .predicates}}
            <a href="https://www.accessdata.fda.gov/scripts/cdrh/cfdocs/cfPMN/pmn.cfm?ID={{.}}">{{ . }}</a>,
            {{ end }}
        </td>
    </tr>
    {{ end }}
    </table>
    {{if .ShowPrev }}
    <a href="/classic/search?query={{.OriginalQuery.Query}}&page={{.PrevPage}}&maxHits={{.MaxResults}}"> <p>Previous Page</p></a>
    {{ end }}
    {{ if .MoreResults }}
    <a href="/classic/search?query={{.OriginalQuery.Query}}&page={{.PrevPage}}&maxHits={{.MaxResults}}"> <p>Next Page</p></a>
    {{ end }}
    <script defer data-domain="fda510k.navan.dev" src="https://plausible.io/js/script.js"></script>
</body>
</html>