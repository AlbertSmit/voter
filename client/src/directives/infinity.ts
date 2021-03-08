export function infinity(node: HTMLElement) {
  const w: HTMLElement | null = node;
  const c: ChildNode | null = node.firstChild;
  const b: ChildNode | null = node.firstChild!.firstChild;

  /**
   * Handle the scrolling event.
   * Two different things going on:
   *
   * 1.   Scroll upwards
   * 2.   Scroll downwards
   */
  const handleScroll = (): void => {
    console.log("w", w.offsetHeight);
    console.log("c", (c as HTMLElement).offsetHeight);
    console.log("b", (b as HTMLElement).offsetHeight);
    // 1.
    if (w.scrollTop < 1) {
      w.scrollTop =
        (c as HTMLElement).offsetHeight - (b as HTMLElement).offsetHeight;
    }

    // 2.
    if (w.scrollTop + w.offsetHeight > (c as HTMLElement).offsetHeight) {
      w.scrollTop = (b as HTMLElement).offsetHeight - w.offsetHeight;
    }
  };

  /**
   * Listen for the
   * scroll event.
   */
  node.addEventListener("scroll", handleScroll);

  /**
   * Set scroll to 1px
   * to trick browser
   * to be able to
   * scroll up.
   */
  w.scrollTop = 1;

  return {
    destroy() {
      node.removeEventListener("scroll", handleScroll);
    },
  };
}
