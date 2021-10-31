import {
    CarouselProvider,
    Slider,
    Slide,
    ButtonBack,
    ButtonNext,
} from "pure-react-carousel";
import "pure-react-carousel/dist/react-carousel.es.css";
import { makeStyles } from "@mui/styles";
import ArrowBackIosIcon from "@mui/icons-material/ArrowBackIos";
import ArrowForwardIosIcon from "@mui/icons-material/ArrowForwardIos";
import { NodeCard } from "..";
const useStyles = makeStyles(() => ({
    "carousel-provider-style": {
        width: "100%",
        display: "flex",
        flexDirection: "row",
    },
    "carousel-button-style": {
        background: "none",
        border: "none",
        outline: "none",
        height: "24px",
        margin: "0",
        alignSelf: "center",
    },
    "carousel-slider-style": {
        width: "100%",
        display: "flex",
    },
    "carousel-slide-style": {
        width: "100%",
        height: "100%",
    },
}));

const Carousel = () => {
    const classes = useStyles();
    return (
        <CarouselProvider
            step={1}
            visibleSlides={4}
            dragStep={1}
            totalSlides={8}
            infinite={true}
            naturalSlideWidth={100}
            naturalSlideHeight={100}
            isIntrinsicHeight={true}
            className={classes["carousel-provider-style"]}
        >
            <ButtonBack
                className={classes["carousel-button-style"]}
                onClick={() => {}}
            >
                <ArrowBackIosIcon name="back" />
            </ButtonBack>
            <Slider className={classes["carousel-slider-style"]}>
                {[1, 2, 3, 4, 5, 6, 7, 8].map(id => (
                    <Slide
                        key={id}
                        index={id}
                        className={classes["carousel-slide-style"]}
                    >
                        <NodeCard isConfigure={true} />
                    </Slide>
                ))}
            </Slider>

            <ButtonNext
                className={classes["carousel-button-style"]}
                onClick={() => {}}
            >
                <ArrowForwardIosIcon name="next" />
            </ButtonNext>
        </CarouselProvider>
    );
};
export default Carousel;
