
res = lib.ParseHtml(input)

output = JSON.stringify({
    "aaa": "bbb",
    "input": input,
    "res": res.FindSingle("head")
})