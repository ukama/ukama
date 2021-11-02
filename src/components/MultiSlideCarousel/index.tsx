import { Box } from "@mui/material";
import { styled } from "@mui/system";
import "@brainhubeu/react-carousel/lib/style.css";
import Carousel, { slidesToShowPlugin } from "@brainhubeu/react-carousel";

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
                    "arrows",
                    {
                        resolve: slidesToShowPlugin,
                        options: {
                            numberOfSlides: numberOfSlides,
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
