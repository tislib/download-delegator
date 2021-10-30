
res = lib.Parser.ParseHtml(input)

output = JSON.stringify({
    "aaa": "bbb",
    "input": input,
    "res": res.FindSingle("p").Text()
})