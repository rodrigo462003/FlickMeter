import { modal } from './modal.js';
const observer = new MutationObserver(() => {
    const codeContainer = document.getElementById('codeContainer');
    if (codeContainer) {
        const digitInputs = document.getElementsByClassName('digitInput');
        for (const input of digitInputs) {
            input.addEventListener('beforeinput', moveFocus.bind(null, input));
        }
    }
});
if (modal)
    observer.observe(modal, { childList: true });
const forwardAndWrite = (input, digit) => {
    const next = input.nextElementSibling;
    next === null || next === void 0 ? void 0 : next.focus();
    if (next)
        next.value = digit;
};
const writeAndForward = (input, digit) => {
    input.value = digit;
    const next = input.nextElementSibling;
    next === null || next === void 0 ? void 0 : next.focus();
};
const backwardsAndDelete = (input) => {
    const prev = input.previousSibling;
    prev === null || prev === void 0 ? void 0 : prev.focus();
    if (prev)
        prev.value = '';
};
const deleteAndBackwards = (input) => {
    input.value = '';
    const prev = input.previousSibling;
    prev === null || prev === void 0 ? void 0 : prev.focus();
};
const paste = (digits) => {
    const digitInputs = document.getElementsByClassName('digitInput');
    if (digitInputs.length !== 6 || digits.length !== 6) {
        return;
    }
    for (let i = 0; i < 6; i++) {
        digitInputs[i].value = digits[i];
    }
    digitInputs[5].focus();
};
const moveFocus = (input, event) => {
    event.preventDefault();
    if (event.inputType === "deleteContentBackward" || event.inputType === "deleteContentForward") {
        if (input.value.length === 0) {
            backwardsAndDelete(input);
            return;
        }
        deleteAndBackwards(input);
        return;
    }
    const text = event.data || "";
    if (event.inputType === "insertFromPaste") {
        if (/^\d{6}$/.test(text)) {
            paste([...text]);
        }
        return;
    }
    if (/^[0-9]$/i.test(text)) {
        if (input.value.length === 0) {
            writeAndForward(input, text);
            return;
        }
        forwardAndWrite(input, text);
    }
};
