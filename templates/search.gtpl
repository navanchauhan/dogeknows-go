<!DOCTYPE html>
<html>
<head>
    <title>{{.GlobalVars.Name}} - Classic Mode</title>
    <script defer data-domain="fda510k.navan.dev" src="https://plausible.io/js/script.js"></script>
</head>
<body>
    <h1>{{.GlobalVars.Name}} - Classic Mode</h1>
    <p>Classic mode makes it easier to scrap results and is intended for powerusers. Not updated often, but will always work without any JS</p>
    <form action="/classic/search" method="GET">
        <input type="text" name="query" placeholder="Search Query" spellcheck="false">
        <input type="submit">
    </form>
    <script defer data-domain="fda510k.navan.dev" src="https://plausible.io/js/script.js"></script>
</body>
</html>