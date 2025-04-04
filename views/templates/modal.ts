const hideForm = (event: Event) => {
    const registerContainer = document.getElementById('registerContainer') as HTMLDivElement | null;
    const codeContainer = document.getElementById("codeContainer") as HTMLDivElement | null;
    const e = event as CustomEvent
    if (e.detail.serverResponse.includes("codeForm")) {
        if (registerContainer) registerContainer.style.display = 'none';
        if (codeContainer) codeContainer.style.display = '';
        return
    }
    if (e.detail.serverResponse.includes("usernameErr")) {
        if (registerContainer) registerContainer.style.display = '';
        if (codeContainer) codeContainer.style.display = 'none';
    }
}

export const modal = document.getElementById('modal') as HTMLDialogElement | null;
modal?.addEventListener('keydown', event => {
    if (!modal.open) {
        return
    }
    if (event.key === 'Escape') {
        animationAndClose(event)
    }

});

modal?.addEventListener('htmx:before-swap', hideForm);

modal?.addEventListener('mousedown', (event) => {
    if (!modal.open) {
        return
    }
    if (event.buttons == 1 && event.target === modal) {
        animationAndClose(event)
    }
});

modal?.addEventListener('htmx:after-settle', (event: any) => {
    event.target?.id === 'modal' && event.detail.xhr.status === 200 ? toggleModal() : null
});

function animationAndClose(event: Event) {
    event.preventDefault();

    const isAnimationRunning = !modal?.classList.contains('closing') && modal?.getAnimations().some(animation => animation.playState === 'running');
    if (isAnimationRunning) {
        modal?.addEventListener('animationend', () => { toggleModal(); }, { once: true });
        return
    }
    toggleModal();
}

const toggleModal = () => {
    if (modal?.open) {
        modal.focus()
        modal.addEventListener('animationend', () => {
            modal.close();
            modal.classList.remove('closing');
        }, { once: true });

        modal.classList.add('closing');
        return
    }

    const input = document.getElementById('form')?.querySelector('input');

    modal?.addEventListener('animationend', function() {
        input?.focus();
    }, { once: true });
    modal?.showModal();
    modal?.focus();
}
