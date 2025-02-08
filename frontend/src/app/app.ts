import "styles/app.css";
import "./vendor";

type Email = {
  image: string;
  title: string;
  description: string;
  needResponse: boolean;
};

const emails: Array<Email> = [
  {
    image: "/images/slide-1.png",
    title: "Docusign Request",
    description:
      "Formalize your agreement with the seller by signing the contract.",
    needResponse: false,
  },
  {
    image: "/images/slide-2.png",
    title: "Condo Closing Preparation",
    description:
      "Here are the documents you need to prepare for your condo closing.",
    needResponse: true,
  },
  {
    image: "/images/slide-3.png",
    title: "Gift Letter Clarification",
    description: "We need more information about the gift letter you provided.",
    needResponse: true,
  },
];

const renderDate = () => {
  const dateElement: HTMLElement = document.getElementById("date");
  const date = new Date();
  const formattedDate = date.toLocaleDateString("en-US");
  const hours = date.getHours();
  const minutes = date.getMinutes();
  const formattedTime = `${hours}:${minutes < 10 ? "0" : ""}${minutes}`;
  dateElement.innerHTML = `<span class="u-color-basic-brighter">${formattedDate}</span> at <span class="u-color-basic-brighter">${formattedTime}</span>`;
};

const renderEmailCarousel = (emails: Array<Email>) => {
  const carouselElement: HTMLElement =
    document.getElementById("emailsCarousel");
  const carouselItemTemplate = document.getElementById("emailCarouselItem");
  let firstCarouselItemActivated = false;
  emails.forEach(function (email: Email) {
    const clone = (<HTMLTemplateElement>carouselItemTemplate).content.cloneNode(
      true
    );

    const captionItem: HTMLElement = (<HTMLElement>clone).querySelector(
      ".carousel-item"
    );
    if (!firstCarouselItemActivated) {
      captionItem.classList.add("active");
      firstCarouselItemActivated = true;
    }

    let image: HTMLElement = (<HTMLElement>clone).querySelector(".item-image");
    image.style.backgroundImage = `url(${email.image})`;

    const captionTitle: HTMLElement = (<HTMLElement>clone).querySelector(
      ".carousel-caption h5"
    );
    const captionTitleText = `${email.title} ${
      email.needResponse
        ? '<span class="badge bg-secondary">Needs Response</span>'
        : ""
    }`;
    captionTitle.innerHTML = captionTitleText;

    const captionDescription = (<HTMLElement>clone).querySelector(
      ".carousel-caption p"
    );
    captionDescription.innerHTML = email.description;

    carouselElement.appendChild(clone);
  });
};

const renderApp = () => {
  renderEmailCarousel(emails);
  renderDate();
};

renderApp();

// re-render every 30 seconds
setInterval(() => {
  renderApp();
}, 30000);
