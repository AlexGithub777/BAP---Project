(function () {
    "use strict";

    // Fetch all the forms we want to apply custom Bootstrap validation styles to
    var forms = document.querySelectorAll(".needs-validation");

    // Loop over them and prevent submission
    Array.prototype.slice.call(forms).forEach(function (form) {
        form.addEventListener(
            "submit",
            function (event) {
                if (!form.checkValidity()) {
                    event.preventDefault();
                    event.stopPropagation();
                } else {
                    // Get the submit button
                    var submitButton = form.querySelector("#submitButton");

                    // Disable the button
                    submitButton.disabled = true;

                    // Show the spinner
                    submitButton
                        .querySelector(".spinner-border")
                        .classList.remove("d-none");

                    // Change the button text
                    submitButton.querySelector(".button-text").textContent =
                        "Sending...";

                    // The form will now submit normally, and the page will redirect
                }

                form.classList.add("was-validated");
            },
            false
        );
    });
})();
