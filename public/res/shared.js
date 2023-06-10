window.onload = function () {
	let modal = document.getElementById("modal")
	function toggleModal() {
		if (modal.hidden) {
			modal.hidden = false
			window.scroll({
				top: 10000,
				behavior: "smooth",
			});
		} else {
			window.scroll({
				top: window.scrollY - modal.clientHeight - 1,
				behavior: "smooth",
			});
			setTimeout(() => {
				modal.hidden = true
			}, 366)
		}
	}

	let settingsButton = document.getElementById("settings")
	settingsButton.addEventListener("click", toggleModal)
	let closeButton = document.getElementById("closebutton")
	closeButton.addEventListener("click", toggleModal)
};