export function infinity(node: Node) {
  let wrapper: HTMLElement | null = document.getElementById("wrapper");
  let content: HTMLElement | null = document.getElementById("content");
  let block: HTMLElement | null = document.getElementById("block");

  /**
   * Handle the scrolling event.
   * Two different things going on:
   *
   * 1.   Scroll upwards
   * 2.   Scroll downwards
   */
  const handleScroll = (): void => {
    // 1.
    if (wrapper!.scrollTop === 0) {
      wrapper!.scrollTop = content!.offsetHeight - block!.offsetHeight;
    }

    // 2.
    if (wrapper!.scrollTop + wrapper!.offsetHeight > content!.offsetHeight) {
      wrapper!.scrollTop = block!.offsetHeight - wrapper!.offsetHeight;
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
  wrapper!.scrollTop = 1;

  return {
    destroy() {
      node.removeEventListener("scroll", handleScroll);
    },
  };
}
