<!DOCTYPE html>
<html>
<body>
    <h1>DogeKnows</h1>
    <form action="/classic/search" method="GET">
        <input type="text" name="query" value="{{.OriginalQuery.Query}}" placeholder="Search Query" spellcheck="false">
        <input type="submit">
    </form>
</body>
</html>