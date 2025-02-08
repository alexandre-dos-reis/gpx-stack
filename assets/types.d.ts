import htmxType from "htmx.org";
import { Alpine } from "alpinejs";

declare global {
  var Alpine: Alpine;
  var htmx: typeof htmxType;
}

declare global {
  interface Window {
    Alpine: Alpine;
    htmx: typeof htmxType;
  }
}
