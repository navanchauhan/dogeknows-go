<!DOCTYPE html>
<html>
<body>
    <h1>{{.GlobalVars.Name}} - Classic Mode</h1>
    <form action="/classic/search" method="GET">
        <input type="text" name="query" placeholder="Search Query" spellcheck="false">
        <input type="submit">
    </form>
</body>
</html>