function render(people) {
    let tbodyHtml = "";
    for (const person of people) {
        tbodyHtml += `<tr>
    <td>${person.id}</td>
    <td>${person.age}</td>
    <td>${person.name}</td>
</tr>`;
    }
    document.getElementById("tbody").innerHTML = tbodyHtml;
}

function refresh() {
    fetch("/api/people").then(function (response) {
        response.json().then(render)
    });
}


function addPerson() {
    const form = document.forms["person"];
    const person = {
        name: form["name"].value,
        age: Number(form["age"].value),
    };

    fetch("/api/people", {method: "POST", body: JSON.stringify(person), headers: {
            'content-type': 'application/json'
        }}).then(refresh);

    form["name"].value = null;
    form["age"].value = null;

    return false;
}

refresh();