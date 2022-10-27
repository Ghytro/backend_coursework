function pollOptionId(index) {
    return "poll-option-" + index;
}

function addOption() {
    let parent = document.getElementById("added-options");
        newInput = document.createElement("input"),
        inputIdx = document.getElementsByClassName("bounding-div").length,
        boundingDiv = document.createElement("div");
    
    boundingDiv.setAttribute("id", pollOptionId(inputIdx));

    newInput.setAttribute("class", "card-input");
    newInput.setAttribute("type", "text");
    newInput.setAttribute("placeholder", "Вариант №"+(inputIdx + 1));
    newInput.setAttribute("name", "options[]");
    boundingDiv.appendChild(newInput);

    let deleteButton = document.createElement("button");
    newInput.setAttribute("value", "X");
    newInput.setAttribute("onclick", "document.getElementById(pollOptionId("+inputIdx+")).remove()");
    newInput.setAttribute("style", "margin-left: 25px;")
    boundingDiv.appendChild(deleteButton);

    parent.appendChild(boundingDiv);
}
