import { Admin, Resource, ListGuesser, EditGuesser } from "react-admin";
import simpleRestProvider from "ra-data-simple-rest";

const dataProvider = simpleRestProvider("http://localhost:3000/admin");

const App = () => (
  <Admin dataProvider={dataProvider}>
    <Resource name="products" list={ListGuesser} edit={EditGuesser} />
  </Admin>
);

export default App;
