function pollOptionId(index) {
    return "poll-option-" + index;
}

function addOption() {
    let parent = document.getElementById("added-options");
        newInput = document.createElement("input"),
        inputIdx = parent.children.length + 1,
        boundingDiv = document.createElement("div");
    
    boundingDiv.setAttribute("id", pollOptionId(inputIdx));
    boundingDiv.setAttribute("class", "bounding-div")

    newInput.setAttribute("class", "card-input");
    newInput.setAttribute("type", "text");
    newInput.setAttribute("placeholder", "Вариант №"+(inputIdx + 1));
    newInput.setAttribute("name", "options[]");
    boundingDiv.appendChild(newInput);

    let deleteButton = document.createElement("input");
    deleteButton.setAttribute("type", "button");
    deleteButton.setAttribute("value", "X");
    deleteButton.setAttribute("onclick", "removePollOption("+inputIdx+")");
    deleteButton.setAttribute("style", "margin-left: 25px;")
    boundingDiv.appendChild(deleteButton);

    parent.appendChild(boundingDiv);
}

function removePollOption(idx) {
    document.getElementById(pollOptionId(idx)).remove();
    let inputs = document.querySelectorAll('#added-options .bounding-div input[type="text"]');
    for (let i = 0; i < inputs.length; i++) {
        inputs[i].setAttribute("placeholder", "Вариант №" + (i + 2));
    }
    let boundingDivs = document.getElementsByClassName("bounding-div");
    for (let i = 0; i < boundingDivs.length; i++) {
        boundingDivs[i].setAttribute("id", pollOptionId(i + 1))
    }
}
