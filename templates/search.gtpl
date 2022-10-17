<!DOCTYPE html>
<html>
<body>
    <h1>{{.GlobalVars.Name}}</h1>
    <form action="/classic/search" method="GET">
        <input type="text" name="query" value="{{.OriginalQuery.Query}}" placeholder="Search Query" spellcheck="false">
        <input type="submit">
    </form>
</body>
</html>