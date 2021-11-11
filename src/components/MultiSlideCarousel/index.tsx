import Carousel, {
    slidesToShowPlugin,
    arrowsPlugin,
} from "@brainhubeu/react-carousel";
import { Box } from "@mui/material";
import { styled } from "@mui/system";
import "@brainhubeu/react-carousel/lib/style.css";
import { ChevronLeft, ChevronRight } from "@mui/icons-material";

const Container = styled(Box)({
    flexFlow: "row",
    display: "inline-grid",
    justifyContent: "center",
    alignContent: "center",
    textAlign: "center",
});

const MultiSlideCarousel = (props: any) => {
    const { children, numberOfSlides = 3 } = props;

    return (
        <Container>
            <Carousel
                plugins={[
                    {
                        resolve: slidesToShowPlugin,
                        options: {
                            numberOfSlides: numberOfSlides,
                        },
                    },
                    {
                        resolve: arrowsPlugin,
                        options: {
                            numberOfSlides: numberOfSlides,
                            arrowLeft: <ChevronLeft />,
                            arrowLeftDisabled: <ChevronLeft />,
                            arrowRight: <ChevronRight />,
                            arrowRightDisabled: <ChevronRight />,
                            addArrowClickHandler: true,
                        },
                    },
                ]}
            >
                {children}
            </Carousel>
        </Container>
    );
};
export default MultiSlideCarousel;
