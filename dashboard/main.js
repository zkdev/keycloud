var exampleEntries = [
    {
        "i" : 1,
        "Username" : "Mark",
        "Url" : "https://www.google.com",
    },{
        "i" : 2,
        "Username" : "Jacob",
        "Url" : "https://www.google.com",
    },{
        "i" : 3,
        "Username" : "Larry",
        "Url" : "https://www.google.com",
    }
]

function clickTab(elem){
    for(let i = 0; i < $("a.nav-link").length; i++){
        $("a.nav-link")[i].classList.remove("active");
    }
    elem.classList.add("active");
}

function addCustomField(){
    $(`<div class="form-row custom-field-row-added" style="margin-bottom: 15px">
                                    <div class="col">
                                        <input type="text" class="form-control" placeholder="Custom Field Name">
                                    </div>
                                    <div class="col">
                                        <input type="text" class="form-control" placeholder="Custom Field content">
                                    </div>
                                    <div class="col">
                                        <input class="form-check-input big-checkbox" type="checkbox">
                                        <label class="form-check-label" style="font-size: x-large;margin-left: 15px;">
                                            Encrypt
                                        </label>
                                    </div>
                                </div>`).insertBefore("#btn-add-field-group");
}

function addTableRow(value) {
    $("#pwTable").prepend(`<tr>
    <th scope="row">${value.i}</th>
    <td>${value.Username}</td>
    <td><a target="_blank" rel="noopener noreferrer" href="${value.Url}">${value.Url}</a></td>
    <td><button type="button" class="btn btn-info"><i class="fa fa-clipboard"></i> Copy to Clipboard</button></td>
    <td><button type="button" class="btn btn-danger"><i class="fa fa-remove"></i></button></td>
    </tr>`);
}

function updateModal(){
    $(".custom-field-row-added").remove();
}

function removeEntry(event) {
    console.log(event);
}

function saveNewEntry() {
    let newEntry = {};
    let formElements = document.forms["newEntryForm"].getElementsByTagName("input");
    newEntry.Url = formElements[0].value;
    exampleEntries.push(newEntry);
    renderTable();
}

function renderTable() {
    exampleEntries.reverse();  // bc of callback
    exampleEntries.forEach(addTableRow);
}

$("document").ready(function() {
    renderTable();
})
