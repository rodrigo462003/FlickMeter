const reviewForm = document.getElementById('newReviewForm') as HTMLFormElement | null;
declare let base: number;

const newReviewButton = document.getElementById('newReviewButton') as HTMLButtonElement | null;
newReviewButton?.addEventListener('click', () => {
    reviewForm?.classList.toggle('h-[340px]');
    reviewForm?.classList.toggle('h-0');
});

const starsHover = (i: number) => {
    const stars = document.getElementsByClassName("buttonStar") as HTMLCollectionOf<HTMLButtonElement>;
    let j = 0
    for (; j < i + 1; j++) {
        (stars[j].firstChild as SVGSVGElement).setAttribute("fill-opacity", "1");
    }
    for (; j < 10; j++) {
        (stars[j].firstChild as SVGSVGElement).setAttribute("fill-opacity", "0");
    }
}

const newBase = (i: number) => {
    base = i + 1;
    resetStars();
}

const resetStars = () => {
    const stars = document.getElementsByClassName("buttonStar") as HTMLCollectionOf<HTMLButtonElement>;
    const starInput = document.getElementById("starInput") as HTMLInputElement | null;
    let j = 0
    if (starInput) starInput.value = base.toString()
    for (; j < base; j++) {
        (stars[j].firstChild as SVGSVGElement).setAttribute("fill-opacity", "1");
    }
    for (; j < 10; j++) {
        (stars[j].firstChild as SVGSVGElement).setAttribute("fill-opacity", "0");
    }
}

reviewForm?.addEventListener('htmx:after-settle', resetStars);


(window as any).isReviewFormVisible = () => {
    return !reviewForm?.classList.contains('h-0');
};
