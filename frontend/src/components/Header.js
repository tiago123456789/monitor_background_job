import React from "react"
import { Navbar, NavbarBrand } from "reactstrap"

export default () => {
    return (
        <div>
            <Navbar color="dark" dark>
                <NavbarBrand href="/" className="mr-auto" style={{ "marginLeft": "10px" }}>
                    Monitor Jobs
          </NavbarBrand>
            </Navbar>
        </div>
    )
}