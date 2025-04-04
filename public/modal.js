const hideForm = (event) => {
    const registerContainer = document.getElementById('registerContainer');
    const codeContainer = document.getElementById("codeContainer");
    const e = event;
    if (e.detail.serverResponse.includes("codeForm")) {
        if (registerContainer)
            registerContainer.style.display = 'none';
        if (codeContainer)
            codeContainer.style.display = '';
        return;
    }
    if (e.detail.serverResponse.includes("usernameErr")) {
        if (registerContainer)
            registerContainer.style.display = '';
        if (codeContainer)
            codeContainer.style.display = 'none';
    }
};
export const modal = document.getElementById('modal');
modal === null || modal === void 0 ? void 0 : modal.addEventListener('keydown', event => {
    if (!modal.open) {
        return;
    }
    if (event.key === 'Escape') {
        animationAndClose(event);
    }
});
modal === null || modal === void 0 ? void 0 : modal.addEventListener('htmx:before-swap', hideForm);
modal === null || modal === void 0 ? void 0 : modal.addEventListener('mousedown', (event) => {
    if (!modal.open) {
        return;
    }
    if (event.buttons == 1 && event.target === modal) {
        animationAndClose(event);
    }
});
modal === null || modal === void 0 ? void 0 : modal.addEventListener('htmx:after-settle', (event) => {
    var _a;
    ((_a = event.target) === null || _a === void 0 ? void 0 : _a.id) === 'modal' && event.detail.xhr.status === 200 ? toggleModal() : null;
});
function animationAndClose(event) {
    event.preventDefault();
    const isAnimationRunning = !(modal === null || modal === void 0 ? void 0 : modal.classList.contains('closing')) && (modal === null || modal === void 0 ? void 0 : modal.getAnimations().some(animation => animation.playState === 'running'));
    if (isAnimationRunning) {
        modal === null || modal === void 0 ? void 0 : modal.addEventListener('animationend', () => { toggleModal(); }, { once: true });
        return;
    }
    toggleModal();
}
const toggleModal = () => {
    var _a;
    if (modal === null || modal === void 0 ? void 0 : modal.open) {
        modal.focus();
        modal.addEventListener('animationend', () => {
            modal.close();
            modal.classList.remove('closing');
        }, { once: true });
        modal.classList.add('closing');
        return;
    }
    const input = (_a = document.getElementById('form')) === null || _a === void 0 ? void 0 : _a.querySelector('input');
    modal === null || modal === void 0 ? void 0 : modal.addEventListener('animationend', function () {
        input === null || input === void 0 ? void 0 : input.focus();
    }, { once: true });
    modal === null || modal === void 0 ? void 0 : modal.showModal();
    modal === null || modal === void 0 ? void 0 : modal.focus();
};
