workspace {

    model {
        user = person "Node"
        softwareSystem = softwareSystem "Bootstrap System" {
            !docs docs
        }

        user -> softwareSystem "Uses"
    }

    views {
        systemContext softwareSystem "Arch" {
            include *
            autoLayout
        }

        theme default
    }

}
