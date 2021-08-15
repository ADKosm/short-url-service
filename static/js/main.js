function isValidURL(url) {
    // see http://urlregex.com/
    var pattern = /((([A-Za-z]{3,9}:(?:\/\/)?)(?:[\-;:&=\+\$,\w]+@)?[A-Za-z0-9\.\-]+|(?:www\.|[\-;:&=\+\$,\w]+@)[A-Za-z0-9\.\-]+)((?:\/[\+~%\/\.\w\-_]*)?\??(?:[\-\+=&;%@\.\w_]*)#?(?:[\.\!\/\\\w]*))?)/;
    return pattern.test(url);
}

function addNewURL(longUrl) {
    var mainInput = document.getElementById("main-input");
    var mainComponent = document.getElementById("main-component");
    var snackbarContainer = document.getElementById('main-snakbar');
    var data;

    if (isValidURL(longUrl)) {
        $.ajax({
            method: "POST",
            url: "api/add",
            dataType: "json",
            contentType: "application/json",
            data: JSON.stringify({url: longUrl})
        }).done(function (msg) {
            if (msg.ShortUrl) {
                var shortUrl = window.location.origin + "/" + msg.ShortUrl;

                mainInput.value = shortUrl;
                mainComponent.classList.add("is-dirty");

                mainInput.focus();
                mainInput.select();

                data = {message: 'Succesfully shorted!'};
            } else {
                data = {message: 'Unknown error!'};
            }
            snackbarContainer.MaterialSnackbar.showSnackbar(data);
        });
    } else {
        mainInput.value = longUrl;
        mainComponent.classList.add("is-dirty");

        data = {message: 'Incorrect url!'};
        snackbarContainer = document.querySelector('#main-snakbar');
        snackbarContainer.MaterialSnackbar.showSnackbar(data);
    }
}

(function main() {
    var mainComponent = document.getElementById("main-component");
    var mainButton = document.getElementById("main-button");
    var mainInput = document.getElementById("main-input");

    mainComponent.addEventListener('paste', function (event) {
        if(event.clipboardData && event.clipboardData.getData) {
            event.preventDefault();
            var url = event.clipboardData.getData('text/plain');
            addNewURL(url);
        }
    });

    mainButton.addEventListener('click', function (event) {
        var url = mainInput.value;
        addNewURL(url);
    });

    mainInput.addEventListener('keypress', function (event) {
        var key = event.which || event.keyCode;
        if(key === 13) { // Enter
            var url = mainInput.value;
            addNewURL(url);
        }
    });
})();