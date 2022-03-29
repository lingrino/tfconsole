const historyKey = "tfconsolehistory"

// things to do when the page loads
document.body.onload = onLoad
function onLoad() {
    addLine()
    fillHistory()
}

// add a new line to the command history
function addLine() {
    // don't add a line if there's no id
    id = document.getElementById("id").innerText || ""
    if (id == "") {
        return
    }

    // don't add a line if it's already been added
    hist = localStorage.getItem(historyKey) || "[]"
    hist = JSON.parse(hist)
    if (hist.length > 0) {
        latestLine = hist[hist.length - 1]
        if (latestLine[0] == id) {
            return
        }
    }

    // add a new line to local storage
    input = document.getElementById("input").innerText
    output = document.getElementById("output").innerText
    hist.push([id, input, output])

    localStorage.setItem(historyKey, JSON.stringify(hist))
}

// get console history as a list of lines
function getLines() {
    hist = localStorage.getItem(historyKey) || "[]"
    return JSON.parse(hist)
}

// populate the command history from local storage
function fillHistory() {
    let html
    for (const line of getLines()) {
        html += "<div id=\""+line[0]+"\" class=\"line\">> "+line[1]+"</br>"+line[2]+"</br></div>"
    }

    const escapeHTMLPolicy = trustedTypes.createPolicy("myEscapePolicy", {
        createHTML: (string) => string // TODO actually filter the lines
    });

    document.getElementById("history").innerHTML = escapeHTMLPolicy.createHTML(html)
}

// things to do before the form is submitted
document.getElementById("form").onsubmit = onSubmit
function onSubmit() {
    setID()
}

// set a unique ID that can be used to identify form submissions
function setID() {
    document.getElementById("formid").value = crypto.randomUUID()
}

// clear all command history from local storage
document.getElementById("clearHistory").onclick = clearHistory
function clearHistory() {
    localStorage.removeItem(historyKey)
    fillHistory()
}
