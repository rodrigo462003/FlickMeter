import { modal } from './modal.js';

const observer = new MutationObserver(() => {
    const codeContainer = document.getElementById('codeContainer') as HTMLFormElement | null;
    if (codeContainer) {
        const digitInputs = document.getElementsByClassName('digitInput') as HTMLCollectionOf<HTMLInputElement>;
        for (const input of digitInputs) {
            input.addEventListener('beforeinput', moveFocus.bind(null, input))
        }
    }
});

if (modal) observer.observe(modal, { childList: true });

const forwardAndWrite = (input: HTMLInputElement, digit: string) => {
    const next = input.nextElementSibling as HTMLInputElement | null
    next?.focus();
    if (next) next.value = digit;
}

const writeAndForward = (input: HTMLInputElement, digit: string) => {
    input.value = digit;
    const next = input.nextElementSibling as HTMLInputElement | null;
    next?.focus();
}

const backwardsAndDelete = (input: HTMLInputElement) => {
    const prev = input.previousSibling as HTMLInputElement | null;
    prev?.focus();
    if (prev) prev.value = '';
}

const deleteAndBackwards = (input: HTMLInputElement) => {
    input.value = '';
    const prev = input.previousSibling as HTMLInputElement | null;
    prev?.focus();
}

const paste = (digits: string[]) => {
    const digitInputs = document.getElementsByClassName('digitInput') as HTMLCollectionOf<HTMLInputElement>;
    if (digitInputs.length !== 6 || digits.length !== 6) {
        return;
    }
    for (let i = 0; i < 6; i++) {
        digitInputs[i].value = digits[i];
    }
    digitInputs[5].focus();
}

const moveFocus = (input: HTMLInputElement, event: InputEvent) => {
    event.preventDefault();

    if (event.inputType === "deleteContentBackward" || event.inputType === "deleteContentForward") {
        if (input.value.length === 0) {
            backwardsAndDelete(input);
            return
        }
        deleteAndBackwards(input);
        return;
    }

    const text = event.data || "";
    if (event.inputType === "insertFromPaste") {
        if (/^\d{6}$/.test(text)) {
            paste([...text])
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
}
